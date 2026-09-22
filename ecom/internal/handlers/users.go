package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/abhinayjangde/ecom/internal/httpx"
	middleware "github.com/abhinayjangde/ecom/internal/middlewares"
	"github.com/abhinayjangde/ecom/internal/utils"
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

	// Validate the request
	if err := req.Validate(); err != nil {
		var verr *ValidationError
		if ok := errors.As(err, &verr); ok {
			uh.logger.ErrorContext(ctx, "Validation error", "request_id", requestId, "err", err)
			httpx.Error(w, http.StatusBadRequest, verr.Error(), httpx.CodeValidationFailed)
			return
		}
	}
	// Check if the email already exists
	var exists bool

	existingUser := uh.db.QueryRowContext(
		ctx, `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`, req.Email)

	err := existingUser.Scan(&exists)

	if err != nil {
		uh.logger.ErrorContext(ctx, "Failed to check existing user", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	if exists {
		httpx.Error(w, http.StatusConflict, "user already exists with this email", httpx.CodeConflict)
		return
	}

	// hash the password
	hashedPassword, err := utils.HashPassword(req.Password)

	if err != nil {
		uh.logger.ErrorContext(ctx, "Failed to hash password", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	// creating the user in the database
	row := uh.db.QueryRowContext(ctx, `
	INSERT INTO users (name, email, password_hash)
	VALUES ($1, $2, $3) RETURNING id, email`, req.Name, req.Email, hashedPassword)

	var out CreateUserResponse
	if err := row.Scan(&out.ID, &out.Email); err != nil {
		uh.logger.ErrorContext(ctx, "Failed to insert user", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	uh.logger.InfoContext(ctx, "User created successfully", "request_id", requestId, "user_id", out.ID)
	httpx.WriteJSON(w, http.StatusCreated, out)
}

func (uh UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, "ok")
}
