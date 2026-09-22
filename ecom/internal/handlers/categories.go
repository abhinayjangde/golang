package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/abhinayjangde/ecom/internal/httpx"
	"github.com/google/uuid"

	middleware "github.com/abhinayjangde/ecom/internal/middlewares"
)

type CategoryHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewCategoryHandler(db *sql.DB, logger *slog.Logger) *CategoryHandler {
	return &CategoryHandler{
		db:     db,
		logger: logger,
	}
}

func (ch *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestID := middleware.RequestIDFromContext(ctx)
	userID := middleware.UserIDFromContext(ctx)

	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ch.logger.ErrorContext(ctx, "failed to decode", "request_id", requestID, "user_id", userID, "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid body", httpx.CodeMalformedJSON)
		return
	}

	if err := req.Validate(); err != nil {

		var verr *ValidationError
		if errors.As(err, &verr) {
			ch.logger.ErrorContext(ctx, "validation error", "request_id", requestID, "err", err)
			httpx.Error(w, http.StatusBadRequest, verr.Error(), httpx.CodeValidationFailed)
			return
		}
	}

	// Process the valid request

	row := ch.db.QueryRowContext(ctx, "INSERT INTO categories (name) VALUES ($1) RETURNING id", req.Name)

	var id uuid.UUID

	if err := row.Scan(&id); err != nil {
		ch.logger.ErrorContext(ctx, "failed to insert category", "request_id", requestID, "user_id", userID, "err", err)
		httpx.Error(w, http.StatusConflict, "category already exists", httpx.CodeConflict)
		return
	}

	// Return the response
	httpx.WriteJSON(w, http.StatusCreated, map[string]any{
		"id": id,
	})

}
