package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/abhinayjangde/olxapi/internal/helpers"
	"github.com/redis/go-redis/v9"
)

type listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}

const (
	listingsCacheKey = "listings"
	cacheTTL         = 10 * time.Second
)

type ListingHanlder struct {
	db    *sql.DB
	redis *redis.Client
}

func NewListingHandler(db *sql.DB, redis *redis.Client) *ListingHanlder {
	return &ListingHanlder{
		db:    db,
		redis: redis,
	}
}

func (lh ListingHanlder) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// check if listings are cached in redis
	redisListings, err := lh.redis.Get(ctx, listingsCacheKey).Result()
	if err == nil {
		var listings []listing
		if err := json.Unmarshal([]byte(redisListings), &listings); err != nil {
			slog.ErrorContext(ctx, "json.Unmarshal error",
				"operation", "listings.list",
				"err", err,
			)
			http.Error(w, "error while deserializing cached listings", http.StatusInternalServerError)
			return
		}
		helpers.WriteJSON(w, http.StatusOK, listings)
		return
	} else if err != redis.Nil {
		slog.WarnContext(ctx, "redis cache unavailable",
			"operation", "listings.list",
			"err", err,
		)
	}

	rows, err := lh.db.QueryContext(ctx,
		`SELECT id, title, description, price, city, created_at
			FROM listings
			ORDER BY created_at DESC
			LIMIT 100`,
	)

	if err != nil {
		slog.ErrorContext(ctx, "database query failed",
			"operation", "listings.list",
			"err", err,
		)
		http.Error(w, "error while fetching listings", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	listings := []listing{}

	for rows.Next() {
		var l listing
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
			slog.ErrorContext(ctx, "row scan failed",
				"operation", "listings.list",
				"err", err,
			)
			http.Error(w, "error while deserializing listing", http.StatusInternalServerError)
			return
		}
		listings = append(listings, l)
	}
	if err := rows.Err(); err != nil {
		slog.ErrorContext(ctx, "rows error",
			"operation", "listings.list",
			"err", err,
		)
		http.Error(w, "error while reading listing", http.StatusInternalServerError)
		return
	}

	// caching the listings in redis for 10 seconds
	jsonListings, err := json.Marshal(listings)
	if err != nil {
		slog.ErrorContext(ctx, "json.Marshal error",
			"operation", "listings.list",
			"err", err,
		)
		http.Error(w, "error while serializing listings", http.StatusInternalServerError)
		return
	}
	err = lh.redis.Set(ctx, listingsCacheKey, jsonListings, cacheTTL).Err()
	if err != nil {
		slog.ErrorContext(ctx, "redis.set error",
			"operation", "listings.list",
			"err", err,
		)
		http.Error(w, "error while saving listings to cache", http.StatusInternalServerError)
		return
	}

	slog.InfoContext(ctx, "listings fetched and cached",
		"operation", "listings.list",
		"count", len(listings),
	)

	helpers.WriteJSON(w, http.StatusOK, listings)

}

func (lh ListingHanlder) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	_, err := lh.db.ExecContext(ctx, `DELETE FROM listings WHERE id = $1`, id)
	if err != nil {
		slog.ErrorContext(ctx, "delete error",
			"operation", "listings.delete",
			"err", err,
		)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	slog.InfoContext(ctx, "listing deleted",
		"operation", "listings.delete",
		"listing_id", id,
	)

	helpers.WriteJSON(w, http.StatusNoContent, nil)
}
