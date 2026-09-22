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

	ch.logger.InfoContext(ctx, "category created", "request_id", requestID, "user_id", userID, "category_id", id, "category_name", req.Name)
	// Return the response
	httpx.WriteJSON(w, http.StatusCreated, map[string]any{
		"id": id,
	})

}

func (ch *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := middleware.RequestIDFromContext(ctx)

	// Fetch all categories
	rows, err := ch.db.QueryContext(ctx, "SELECT id, name FROM categories")
	if err != nil {
		ch.logger.ErrorContext(ctx, "Failed to fetch categories", "request_id", requestID, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}
	defer rows.Close()

	var categories []map[string]any
	for rows.Next() {
		var id uuid.UUID
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			ch.logger.ErrorContext(ctx, "Failed to scan category", "request_id", requestID, "err", err)
			httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
			return
		}
		categories = append(categories, map[string]any{
			"id":   id,
			"name": name,
		})
	}

	if err := rows.Err(); err != nil {
		ch.logger.ErrorContext(ctx, "Error occurred while iterating over categories", "request_id", requestID, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	ch.logger.InfoContext(ctx, "categories fetched", "request_id", requestID, "count", len(categories))
	httpx.WriteJSON(w, http.StatusOK, categories)
}

func (ch *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	requestID := middleware.RequestIDFromContext(ctx)
	userID := middleware.UserIDFromContext(ctx)
	if id == "" {
		ch.logger.ErrorContext(ctx, "Missing category ID", "request_id", requestID, "user_id", userID)
		httpx.Error(w, http.StatusBadRequest, "missing category ID", httpx.CodeInvalidID)
		return
	}

	// Delete the category
	_, err := ch.db.ExecContext(ctx, "DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		ch.logger.ErrorContext(ctx, "Failed to delete category", "request_id", requestID, "category_id", id, "user_id", userID, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	ch.logger.InfoContext(ctx, "category deleted", "request_id", requestID, "category_id", id, "user_id", userID)
	w.WriteHeader(http.StatusNoContent)

}
