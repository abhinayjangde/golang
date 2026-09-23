package models

import "time"

type Listings struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Price       int64          `json:"price"`
	City        string         `json:"city"`
	CreatedAt   time.Time      `json:"created_at"`
	Images      []ListingImage `json:"images"`
}

type ListingImage struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	ContentType string `json:"content_type,omitempty"`
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
}
