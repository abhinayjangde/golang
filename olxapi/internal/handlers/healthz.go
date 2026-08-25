package handlers

import (
	"net/http"

	"github.com/abhinayjangde/olxapi/internal/helpers"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJSON(w, http.StatusOK, map[string]any{
		"status": "okay",
	})
}
