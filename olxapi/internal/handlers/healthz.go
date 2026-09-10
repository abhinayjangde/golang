package handlers

import (
	"net/http"

	"github.com/abhinayjangde/olxapi/internal/httpx"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"status": "okay",
	})
}
