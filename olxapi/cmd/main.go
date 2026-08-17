package main

import (
	"log"
	"net/http"
	"time"

	"github.com/abhinayjangde/olxapi/internal/config"
)

func main() {
	cfg := config.MustLoad()
	mux := http.NewServeMux()

	log.Printf("starting server on port %s", cfg.Port)

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(`{"status":"ok"}`))
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("error while starting server %v", err)
	}
}
