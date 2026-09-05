package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/abhinayjangde/olxapi/internal/helpers"
	"github.com/abhinayjangde/olxapi/internal/httpx"
	middleware "github.com/abhinayjangde/olxapi/internal/middlewares"
	"github.com/redis/go-redis/v9"
)

type listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}

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

func (lh ListingHanlder) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)

	// check if listings are cached in redis
	redisListings, err := lh.redis.Get(ctx, listingsCacheKey).Result()
	if err == nil {
		var listings []listing
		if err := json.Unmarshal([]byte(redisListings), &listings); err != nil {
			lh.logger.ErrorContext(ctx, "json.Unmarshal error",
				"operation", "listings.list",
				"err", err,
				"request_id", requestId,
			)
			httpx.Error(w, http.StatusInternalServerError, "error while deserializing cached listings", httpx.CodeInternalError)
			return
		}
		helpers.WriteJSON(w, http.StatusOK, listings)
		return
	} else if err != redis.Nil {
		lh.logger.WarnContext(ctx, "redis cache unavailable",
			"operation", "listings.list",
			"err", err,
			"request_id", requestId,
		)
	}

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

	listings := []listing{}

	for rows.Next() {
		var l listing
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

	helpers.WriteJSON(w, http.StatusOK, listings)

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

	helpers.WriteJSON(w, http.StatusNoContent, nil)
}

func (lh ListingHanlder) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestId := middleware.RequestIDFromContext(ctx)

	var req listing
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lh.logger.ErrorContext(ctx, "failed to decode", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid body", httpx.CodeMalformedJSON)
		return
	}

	// validating req.body
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	req.City = strings.TrimSpace(req.City)

	if req.Title == "" {
		httpx.Error(w, http.StatusBadRequest, "title is required", httpx.CodeMalformedJSON)
		return
	}
	if req.Price <= 0 {
		httpx.Error(w, http.StatusBadRequest, "price must be greater than 0", httpx.CodeMalformedJSON)
		return
	}

	row := lh.db.QueryRowContext(ctx,
		`INSERT INTO listings (title, description, price, city)
			VALUES ($1, $2, $3, $4) RETURNING id`,
		req.Title,
		req.Description,
		req.Price,
		req.City,
	)

	var id string
	if err := row.Scan(&id); err != nil {
		lh.logger.ErrorContext(ctx, "failed to insert", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	if err := lh.invalidateListingsCache(ctx); err != nil && !errors.Is(err, redis.Nil) {
		lh.logger.ErrorContext(ctx, "redis cache invalidation failed",
			"request_id", requestId,
			"err", err,
		)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	lh.logger.InfoContext(ctx, "listing created", "request_id", requestId, "listing_id", id)

	helpers.WriteJSON(w, http.StatusCreated, map[string]string{"status": "created", "id": id})
}
