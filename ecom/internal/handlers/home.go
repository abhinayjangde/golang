package handlers

import (
	"net/http"

	"github.com/abhinayjangde/ecom/internal/httpx"
	middleware "github.com/abhinayjangde/ecom/internal/middlewares"
)

func Home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"github":      "https://github.com/abhinayjangde",
		"source_code": "https://github.com/abhinayjangde/golang/tree/main/ecom",
		"owner":       "abhinayjangde@gmail.com",
		"request_id":  requestId,
	})
}
