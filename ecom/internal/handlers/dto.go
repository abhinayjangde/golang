package handlers

import (
	"fmt"
	"strings"
	"time"
)

// Listing DTOs
type CreateListingRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	City        string `json:"city"`
	CategoryID  string `json:"category_id"`
}

type ValidationError struct {
	Field string
	Msg   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Msg)
}

func (req CreateListingRequest) Validate() error {
	if strings.TrimSpace(req.Title) == "" {
		return &ValidationError{
			Field: "title",
			Msg:   "must not be empty",
		}
	}
	if req.Price <= 0 {
		return &ValidationError{
			Field: "price",
			Msg:   "price must be greater than 0",
		}
	}
	return nil
}

type CreateListingResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

// User DTOs

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (req CreateUserRequest) Validate() error {
	if strings.TrimSpace(req.Name) == "" {
		return &ValidationError{
			Field: "name",
			Msg:   "must not be empty",
		}
	}
	if strings.TrimSpace(req.Email) == "" {
		return &ValidationError{
			Field: "email",
			Msg:   "must not be empty",
		}
	}
	if strings.TrimSpace(req.Password) == "" {
		return &ValidationError{
			Field: "password",
			Msg:   "must not be empty",
		}
	}
	return nil
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req LoginUserRequest) Validate() error {
	if strings.TrimSpace(req.Email) == "" {
		return &ValidationError{
			Field: "email",
			Msg:   "must not be empty",
		}
	}
	if strings.TrimSpace(req.Password) == "" {
		return &ValidationError{
			Field: "password",
			Msg:   "must not be empty",
		}
	}
	return nil
}

type LoginUserResponse struct {
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// profile DTOs
type GetProfileResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

// category DTOs
type CreateCategoryRequest struct {
	Name string `json:"name"`
}

func (req CreateCategoryRequest) Validate() error {
	if strings.TrimSpace(req.Name) == "" {
		return &ValidationError{
			Field: "name",
			Msg:   "must not be empty",
		}
	}
	return nil
}

// image DTOs
var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

type CreateImageRequest struct {
	ContentType string `json:"content_type"`
}

func (req CreateImageRequest) Validate() error {
	if strings.TrimSpace(req.ContentType) == "" {
		return &ValidationError{
			Field: "content_type",
			Msg:   "must not be empty",
		}
	}
	if allowedImageTypes[req.ContentType] == "" {
		return &ValidationError{
			Field: "content_type",
			Msg:   "must be one of: image/jpeg, image/png, image/webp",
		}
	}
	return nil
}

type CreateImageResponse struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	ObjectKey string `json:"object_key"`
	UploadURL string `json:"upload_url"`
	ExpiresIn int64  `json:"expires_in"`
}

type ImageResponse struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	URL         string `json:"url"`
	ContentType string `json:"content_type,omitempty"`
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
}
