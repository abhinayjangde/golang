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
	redis, err := db.NewRedisClient(cfg.RedisUrl)
	if err != nil {
		log.Fatalf("main.db.redis: %v", err)
	}
	defer redis.Close()

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

	log.Println("redis is ready to use!")
	log.Println("postgres database connected")
	log.Printf("starting server on port %s", cfg.Port)

	lh := handlers.NewListingHandler(db, redis) // listing handler

	wrappedMux := c.Handler(mux)
	mux.HandleFunc("GET /", handlers.Home)
	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.HandleFunc("GET /listings", lh.List)
	mux.HandleFunc("DELETE /listings/{id}", lh.Delete)

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
