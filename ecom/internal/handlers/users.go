package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/abhinayjangde/ecom/internal/config"
	"github.com/abhinayjangde/ecom/internal/httpx"
	middleware "github.com/abhinayjangde/ecom/internal/middlewares"
	"github.com/abhinayjangde/ecom/internal/utils"
)

type UserHandler struct {
	db     *sql.DB
	logger *slog.Logger
	cfg    config.Config
}

func NewUserHandler(db *sql.DB, logger *slog.Logger, cfg config.Config) *UserHandler {
	return &UserHandler{
		db:     db,
		logger: logger,
		cfg:    cfg,
	}
}

func (uh UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := middleware.RequestIDFromContext(ctx)
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		uh.logger.ErrorContext(ctx, "Failed to decode request body", "request_id", requestID, "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid body", httpx.CodeMalformedJSON)
		return
	}

	// Validate the request
	if err := req.Validate(); err != nil {
		var verr *ValidationError
		if ok := errors.As(err, &verr); ok {
			uh.logger.ErrorContext(ctx, "Validation error", "request_id", requestID, "err", err)
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
		uh.logger.ErrorContext(ctx, "Failed to check existing user", "request_id", requestID, "err", err)
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
		uh.logger.ErrorContext(ctx, "Failed to hash password", "request_id", requestID, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	// creating the user in the database
	row := uh.db.QueryRowContext(ctx, `
	INSERT INTO users (name, email, password_hash)
	VALUES ($1, $2, $3) RETURNING id, email`, req.Name, req.Email, hashedPassword)

	var out CreateUserResponse
	if err := row.Scan(&out.ID, &out.Email); err != nil {
		uh.logger.ErrorContext(ctx, "Failed to insert user", "request_id", requestID, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	uh.logger.InfoContext(ctx, "User created successfully", "request_id", requestID, "user_id", out.ID)
	httpx.WriteJSON(w, http.StatusCreated, out)
}

func (uh UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := middleware.RequestIDFromContext(ctx)

	// reading request body
	var req LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		uh.logger.ErrorContext(ctx, "Failed to decode request body", "request_id", requestID, "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid body", httpx.CodeMalformedJSON)
		return
	}

	// validating req data
	if err := req.Validate(); err != nil {
		var verr *ValidationError
		if ok := errors.As(err, &verr); ok {
			uh.logger.ErrorContext(ctx, "Validation error", "request_id", requestID, "err", err)
			httpx.Error(w, http.StatusBadRequest, verr.Error(), httpx.CodeValidationFailed)
			return
		}
	}

	// Check if the user exists
	var userID, passwordHash string
	row := uh.db.QueryRowContext(ctx, `SELECT id, password_hash FROM users WHERE email = $1`, req.Email)

	err := row.Scan(&userID, &passwordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			uh.logger.ErrorContext(ctx, "User not found", "request_id", requestID, "email", req.Email)
			httpx.Error(w, http.StatusNotFound, "user not found", httpx.CodeNotFound)
			return
		}
		uh.logger.ErrorContext(ctx, "Failed to query user", "request_id", requestID, "email", req.Email, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	// Verify the password
	if !utils.CheckPasswordHash(req.Password, passwordHash) {
		uh.logger.ErrorContext(ctx, "Invalid password", "request_id", requestID, "email", req.Email)
		httpx.Error(w, http.StatusUnauthorized, "invalid credentials", httpx.CodeUnauthenticated)
		return
	}

	// TODO: If the password is correct, you can generate a token or session here (not implemented in this snippet)
	accessToken, err := utils.GenerateAccessToken(
		userID,
		req.Email,
		uh.cfg.JWTSecret,
	)

	if err != nil {
		uh.logger.ErrorContext(ctx, "error while generating access token", "request_id", requestID, "email", req.Email)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		uh.logger.ErrorContext(ctx, "error while generating refresh token", "request_id", requestID, "email", req.Email)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	var out LoginUserResponse

	out = LoginUserResponse{
		Message:      "User loggged in successfully",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	uh.logger.InfoContext(ctx, "User logged in successfully", "request_id", requestID, "user_id", userID)
	httpx.WriteJSON(w, http.StatusOK, out)
}

func (uh UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := middleware.RequestIDFromContext(ctx)

	userID := middleware.UserIDFromContext(ctx)
	email := middleware.EmailFromContext(ctx)

	out := GetProfileResponse{
		UserID: userID,
		Email:  email,
	}

	uh.logger.InfoContext(ctx, "User profile retrieved successfully", "request_id", requestID, "user_id", userID)
	httpx.WriteJSON(w, http.StatusOK, out)
}
