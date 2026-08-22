package main

import (
	"log"
	"net/http"
	"time"

	"github.com/abhinayjangde/olxapi/internal/config"
	"github.com/abhinayjangde/olxapi/internal/db"
	"github.com/abhinayjangde/olxapi/internal/handlers"
)

func main() {
	cfg := config.MustLoad()
	_, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("main.db.connect: %v", err)
	}
	mux := http.NewServeMux()

	log.Println("database connected")
	log.Printf("starting server on port %s", cfg.Port)

	mux.HandleFunc("GET /healthz", handlers.Healthz)

	srv := http.Server{
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
