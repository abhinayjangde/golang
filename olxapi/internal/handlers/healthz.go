package handlers

import "net/http"

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(`{"status":"okay"}`))
}
