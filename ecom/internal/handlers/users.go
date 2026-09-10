package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/abhinayjangde/ecom/internal/httpx"
	middleware "github.com/abhinayjangde/ecom/internal/middlewares"
)

type UserHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewUserHandler(db *sql.DB, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		db:     db,
		logger: logger,
	}
}

func (uh UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		uh.logger.ErrorContext(ctx, "Failed to decode request body", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid body", httpx.CodeMalformedJSON)
		return
	}

	// TODO: Validate the request
	// TODO: Check if the email already exists
	// TODO: Hash the password
	row := uh.db.QueryRowContext(ctx, `
	INSERT INTO users (name, email, password_hash)
	VALUES ($1, $2, $3) RETURNING id, email`, req.Name, req.Email, req.Password)

	var out CreateUserResponse
	if err := row.Scan(&out.ID, &out.Email); err != nil {
		uh.logger.ErrorContext(ctx, "Failed to insert user", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	uh.logger.InfoContext(ctx, "User created successfully", "request_id", requestId, "user_id", out.ID)
	httpx.WriteJSON(w, http.StatusCreated, out)
}
