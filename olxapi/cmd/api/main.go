package main

import (
	"log"
	"net/http"
	"time"

	"github.com/abhinayjangde/olxapi/internal/config"
	"github.com/abhinayjangde/olxapi/internal/db"
	"github.com/abhinayjangde/olxapi/internal/handlers"
	"github.com/rs/cors"
)

func main() {
	cfg := config.MustLoad()
	redis, err := db.Redis(cfg.RedisUrl)
	if err != nil {
		log.Fatalf("main.db.redis: %v", err)
	}
	db, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("main.db.connect: %v", err)
	}
	mux := http.NewServeMux()
	// wrappedMux := middlewares.CorsMiddleware(mux) // for fixing cors policy
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:5500"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            false, // Set to true to print debugging info to stderr
	})

	wrappedMux := c.Handler(mux)

	log.Println("database connected")
	log.Printf("starting server on port %s", cfg.Port)

	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.HandleFunc("GET /listings", handlers.List(db, redis))

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      wrappedMux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("error while starting server %v", err)
	}
}
