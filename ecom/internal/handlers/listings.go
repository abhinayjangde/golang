package handlers

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/abhinayjangde/ecom/internal/httpx"
	middleware "github.com/abhinayjangde/ecom/internal/middlewares"
	"github.com/abhinayjangde/ecom/internal/models"
	"github.com/redis/go-redis/v9"
)

const (
	listingsCacheKey = "listings"
	cacheTTL         = 20 * time.Second
)

type ListingHanlder struct {
	db     *sql.DB
	redis  *redis.Client
	logger *slog.Logger
}

func NewListingHandler(db *sql.DB, redis *redis.Client, logger *slog.Logger) *ListingHanlder {
	return &ListingHanlder{
		db:     db,
		redis:  redis,
		logger: logger,
	}
}

func (lh ListingHanlder) invalidateListingsCache(ctx context.Context) error {
	return lh.redis.Del(ctx, listingsCacheKey).Err()
}

// Etag hash generation
func generateETag(listings []models.Listings) (string, error) {
	data, err := json.Marshal(listings)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return `"` + hex.EncodeToString(hash[:]) + `"`, nil
}

func writeListingsResponse(w http.ResponseWriter, r *http.Request, listings []models.Listings) error {
	etag, err := generateETag(listings)
	if err != nil {
		return err
	}

	w.Header().Set("ETag", etag)

	if r.Header.Get("If-None-Match") == etag {
		w.WriteHeader(http.StatusNotModified)
		return nil
	}

	httpx.WriteJSON(w, http.StatusOK, listings)
	return nil
}

func (lh ListingHanlder) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)

	// check if listings are cached in redis
	redisListings, err := lh.redis.Get(ctx, listingsCacheKey).Result()
	if err == nil {
		var listings []models.Listings
		if err := json.Unmarshal([]byte(redisListings), &listings); err != nil {
			lh.logger.ErrorContext(ctx, "json.Unmarshal error",
				"operation", "listings.list",
				"err", err,
				"request_id", requestId,
			)
			httpx.Error(w, http.StatusInternalServerError, "error while deserializing cached listings", httpx.CodeInternalError)
			return
		}

		if err := writeListingsResponse(w, r, listings); err != nil {
			lh.logger.ErrorContext(ctx, "failed to write listings response",
				"operation", "listings.list",
				"err", err,
				"request_id", requestId,
			)
			httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
			return
		}
	} else if err != redis.Nil {
		lh.logger.WarnContext(ctx, "redis cache unavailable",
			"operation", "listings.list",
			"err", err,
			"request_id", requestId,
		)
	}

	// if not cached, fetch from database
	rows, err := lh.db.QueryContext(ctx,
		`SELECT id, title, description, price, city, created_at
			FROM listings
			ORDER BY created_at DESC
			LIMIT 100`,
	)

	if err != nil {
		lh.logger.ErrorContext(ctx, "database query failed",
			"operation", "listings.list",
			"err", err,
			"request_id", requestId,
		)

		httpx.Error(w, http.StatusInternalServerError, "error while fetching listings", httpx.CodeInternalError)
		return
	}

	defer rows.Close()

	listings := []models.Listings{}

	for rows.Next() {
		var l models.Listings
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
			lh.logger.ErrorContext(ctx, "row scan failed",
				"operation", "listings.list",
				"err", err,
				"request_id", requestId,
			)
			httpx.Error(w, http.StatusInternalServerError, "error while deserializing listing", httpx.CodeInternalError)
			return
		}
		listings = append(listings, l)
	}
	if err := rows.Err(); err != nil {
		lh.logger.ErrorContext(ctx, "rows error",
			"operation", "listings.list",
			"err", err,
			"request_id", requestId,
		)
		httpx.Error(w, http.StatusInternalServerError, "error while reading listing", httpx.CodeInternalError)
		return
	}

	// caching the listings in redis for 10 seconds
	jsonListings, err := json.Marshal(listings)
	if err != nil {
		lh.logger.ErrorContext(ctx, "json.Marshal error",
			"operation", "listings.list",
			"err", err,
			"request_id", requestId,
		)
		httpx.Error(w, http.StatusInternalServerError, "error while serializing listings", httpx.CodeInternalError)
		return
	}
	err = lh.redis.Set(ctx, listingsCacheKey, jsonListings, cacheTTL).Err()
	if err != nil {
		lh.logger.ErrorContext(ctx, "redis.set error",
			"operation", "listings.list",
			"err", err,
			"request_id", requestId,
		)
		httpx.Error(w, http.StatusInternalServerError, "error while saving listings to cache", httpx.CodeInternalError)
		return
	}

	lh.logger.InfoContext(ctx, "listings fetched and cached",
		"operation", "listings.list",
		"count", len(listings),
		"request_id", requestId,
	)

	if err := writeListingsResponse(w, r, listings); err != nil {
		lh.logger.ErrorContext(ctx, "failed to write listings response",
			"operation", "listings.list",
			"err", err,
			"request_id", requestId,
		)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

}

func (lh ListingHanlder) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	_, err := lh.db.ExecContext(ctx, `DELETE FROM listings WHERE id = $1`, id)
	if err != nil {
		lh.logger.ErrorContext(ctx, "delete error",
			"operation", "listings.delete",
			"err", err,
		)
		httpx.Error(w, http.StatusInternalServerError, "error while deleting a listing", httpx.CodeInternalError)
		return
	}

	if err := lh.invalidateListingsCache(ctx); err != nil && !errors.Is(err, redis.Nil) {
		lh.logger.ErrorContext(ctx, "redis cache invalidation failed",
			"operation", "listings.delete",
			"err", err,
		)
		httpx.Error(w, http.StatusInternalServerError, "redis cache invalidation failed", httpx.CodeInternalError)
		return
	}

	lh.logger.InfoContext(ctx, "listing deleted",
		"operation", "listings.delete",
		"listing_id", id,
	)

	httpx.WriteJSON(w, http.StatusNoContent, nil)
}

func (lh ListingHanlder) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestId := middleware.RequestIDFromContext(ctx)

	var req CreateListingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lh.logger.ErrorContext(ctx, "failed to decode", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid body", httpx.CodeMalformedJSON)
		return
	}

	if err := req.Validate(); err != nil {
		// verr := &ValidationError{}
		var verr *ValidationError
		errors.As(err, &verr)
		httpx.ValidationError(w, http.StatusUnprocessableEntity, err.Error(), httpx.CodeValidationFailed, verr.Field)
		return
	}

	row := lh.db.QueryRowContext(ctx,
		`INSERT INTO listings (title, description, price, city)
			VALUES ($1, $2, $3, $4) RETURNING id, title, created_at`,
		req.Title,
		req.Description,
		req.Price,
		req.City,
	)

	var out CreateListingResponse
	if err := row.Scan(&out.ID, &out.Title, &out.CreatedAt); err != nil {
		lh.logger.ErrorContext(ctx, "failed to insert", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	// invalidate redis cache
	if err := lh.invalidateListingsCache(ctx); err != nil && !errors.Is(err, redis.Nil) {
		lh.logger.ErrorContext(ctx, "redis cache invalidation failed",
			"request_id", requestId,
			"err", err,
		)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	lh.logger.InfoContext(ctx, "listing created", "request_id", requestId, "listing_id", out.ID)

	httpx.WriteJSON(w, http.StatusCreated, out)
}
