package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
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
			log.Printf("json.Unmarshal cached listings: %v", err)
			http.Error(w, "error while deserializing cached listings", http.StatusInternalServerError)
			return
		}
		helpers.WriteJSON(w, http.StatusOK, listings)
		return
	}

	rows, err := lh.db.Query(`SELECT id, title, description, price, city, created_at
			FROM listings
			ORDER BY created_at DESC
			LIMIT 100`)

	if err != nil {
		log.Printf("db.query: %v", err)
		http.Error(w, "error while fetching listings", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	listings := []listing{}

	for rows.Next() {
		var l listing
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
			log.Printf("rows.Scan: %v", err)
			http.Error(w, "error while deserializing listing", http.StatusInternalServerError)
			return
		}
		listings = append(listings, l)
	}
	if err := rows.Err(); err != nil {
		log.Printf("rows.err: %v", err)
		http.Error(w, "error while reading listing", http.StatusInternalServerError)
		return
	}

	// caching the listings in redis for 10 seconds
	jsonListings, err := json.Marshal(listings)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
		http.Error(w, "error while serializing listings", http.StatusInternalServerError)
		return
	}
	err = lh.redis.Set(ctx, listingsCacheKey, jsonListings, cacheTTL).Err()
	if err != nil {
		log.Printf("redis.set: %v", err)
		http.Error(w, "error while saving listings to cache", http.StatusInternalServerError)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, listings)

}

func (lh ListingHanlder) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	_, err := lh.db.Exec(`DELETE FROM listings WHERE id = $1`, id)
	if err != nil {
		log.Printf("delete: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusNoContent, nil)
}
