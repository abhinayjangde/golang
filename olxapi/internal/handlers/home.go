package handlers

import (
	"net/http"

	"github.com/abhinayjangde/olxapi/internal/helpers"
)

func Home(w http.ResponseWriter, r *http.Request) {

	helpers.WriteJSON(w, http.StatusOK, map[string]any{
		"owner":      "abhinayjangde@gmail.com",
		"github":     "https://github.com/abhinayjangde",
		"request_id": r.Context().Value("requestCtxId"),
	})
}
