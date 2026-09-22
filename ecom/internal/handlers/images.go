package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/abhinayjangde/ecom/internal/config"
	"github.com/abhinayjangde/ecom/internal/httpx"
	"github.com/abhinayjangde/ecom/internal/lib"
	middleware "github.com/abhinayjangde/ecom/internal/middlewares"
	"github.com/google/uuid"
)

type ImageHandler struct {
	db     *sql.DB
	s3     *lib.S3Client
	cfg    config.Config
	logger *slog.Logger
}

func NewImageHandler(db *sql.DB, s3c *lib.S3Client, cfg config.Config, logger *slog.Logger) *ImageHandler {
	return &ImageHandler{db: db, s3: s3c, cfg: cfg, logger: logger}
}

func (ih ImageHandler) publicURL(key string) string {
	return fmt.Sprintf(
		"https://%s.s3.%s.amazonaws.com/%s",
		ih.cfg.S3Bucket, ih.cfg.AWSRegion, strings.TrimPrefix(key, "/"),
	)
}

// listIsOwnedBy helper (used by Create)
func (ih ImageHandler) listIsOwnedBy(ctx context.Context, listingID, userID string) bool {
	var one int
	err := ih.db.QueryRowContext(ctx,
		`SELECT 1 FROM listings WHERE id = $1 AND user_id = $2`, listingID, userID,
	).Scan(&one)
	return err == nil
}

func (ih ImageHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := middleware.RequestIDFromContext(ctx)
	userID := middleware.UserIDFromContext(ctx)
	listingID := r.PathValue("id")

	var req CreateImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ih.logger.ErrorContext(ctx, "failed to decode", "request_id", requestID, "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid body", httpx.CodeMalformedJSON)
		return
	}
	if err := req.Validate(); err != nil {
		var verr *ValidationError
		if errors.As(err, &verr) {
			ih.logger.ErrorContext(ctx, "validation error", "request_id", requestID, "err", err)
			httpx.ValidationError(w, http.StatusUnprocessableEntity, verr.Error(), httpx.CodeValidationFailed, verr.Field)
			return
		}
	}

	if !ih.listIsOwnedBy(ctx, listingID, userID) {
		httpx.Error(w, http.StatusNotFound, "listing not found or you are not its owner", httpx.CodeNotFound)
		return
	}

	filename := uuid.NewString() + allowedImageTypes[req.ContentType]
	objectKey := ih.s3.ObjectKey(listingID, filename)

	var imageID string
	err := ih.db.QueryRowContext(ctx,
		`INSERT INTO images (listing_id, object_key, mime, status)
		 VALUES ($1, $2, $3, 'pending') RETURNING id`,
		listingID, objectKey, req.ContentType,
	).Scan(&imageID)
	if err != nil {
		ih.logger.ErrorContext(ctx, "failed to insert image row", "request_id", requestID, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	uploadURL, err := ih.s3.PresignedPutURL(ctx, objectKey, ih.cfg.PresignedTTL)
	if err != nil {
		ih.logger.ErrorContext(ctx, "failed to presign upload url", "request_id", requestID, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	ih.logger.InfoContext(ctx, "image upload url issued",
		"request_id", requestID, "listing_id", listingID, "image_id", imageID)

	httpx.WriteJSON(w, http.StatusCreated, CreateImageResponse{
		ID:        imageID,
		Status:    "pending",
		ObjectKey: objectKey,
		UploadURL: uploadURL,
		ExpiresIn: int64(ih.cfg.PresignedTTL.Seconds()),
	})
}

func (ih ImageHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := middleware.RequestIDFromContext(ctx)
	userID := middleware.UserIDFromContext(ctx)
	imageID := r.PathValue("id")

	tx, err := ih.db.BeginTx(ctx, nil)
	if err != nil {
		ih.logger.ErrorContext(ctx, "failed to begin tx", "request_id", requestID, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx,
		`UPDATE images SET status = 'uploaded'
		 WHERE id = $1
		   AND listing_id IN (SELECT id FROM listings WHERE user_id = $2)`,
		imageID, userID,
	)
	if err != nil {
		ih.logger.ErrorContext(ctx, "failed to confirm upload", "request_id", requestID, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		httpx.Error(w, http.StatusNotFound, "image not found or you are not its owner", httpx.CodeNotFound)
		return
	}

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO image_jobs (image_id, status) VALUES ($1, 'pending')
		 ON CONFLICT (image_id) DO NOTHING`, imageID,
	); err != nil {
		ih.logger.ErrorContext(ctx, "failed to enqueue image job", "request_id", requestID, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	if err := tx.Commit(); err != nil {
		ih.logger.ErrorContext(ctx, "failed to commit tx", "request_id", requestID, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	ih.logger.InfoContext(ctx, "image upload confirmed, job queued", "request_id", requestID, "image_id", imageID)
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"status": "queued"})
}

// Delete: remove the object(s) from S3 and the row.
func (ih ImageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := middleware.RequestIDFromContext(ctx)
	userID := middleware.UserIDFromContext(ctx)
	imageID := r.PathValue("id")

	var objectKey, processedKey string
	err := ih.db.QueryRowContext(ctx,
		`SELECT object_key, COALESCE(processed_key, '')
		 FROM images
		 WHERE id = $1
		   AND listing_id IN (SELECT id FROM listings WHERE user_id = $2)`,
		imageID, userID,
	).Scan(&objectKey, &processedKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpx.Error(w, http.StatusNotFound, "image not found or you are not its owner", httpx.CodeNotFound)
			return
		}
		ih.logger.ErrorContext(ctx, "failed to fetch image", "request_id", requestID, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	if err := ih.s3.DeleteObject(ctx, objectKey); err != nil {
		ih.logger.ErrorContext(ctx, "failed to delete original object", "request_id", requestID, "key", objectKey, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}
	if processedKey != "" && processedKey != objectKey {
		if err := ih.s3.DeleteObject(ctx, processedKey); err != nil {
			ih.logger.ErrorContext(ctx, "failed to delete processed object", "request_id", requestID, "key", processedKey, "err", err)
			httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
			return
		}
	}

	if _, err := ih.db.ExecContext(ctx, `DELETE FROM images WHERE id = $1`, imageID); err != nil {
		ih.logger.ErrorContext(ctx, "failed to delete image row", "request_id", requestID, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	ih.logger.InfoContext(ctx, "image deleted", "request_id", requestID, "image_id", imageID)
	w.WriteHeader(http.StatusNoContent)
}
