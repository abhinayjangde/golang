package handlers

import (
	"net/http"

	"github.com/abhinayjangde/ecom/internal/httpx"
)

func Home(w http.ResponseWriter, r *http.Request) {

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"owner":      "abhinayjangde@gmail.com",
		"github":     "https://github.com/abhinayjangde",
		"request_id": r.Context().Value("requestCtxId"),
	})
}
