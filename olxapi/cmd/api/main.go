package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/abhinayjangde/olxapi/internal/config"
	"github.com/abhinayjangde/olxapi/internal/db"
	"github.com/abhinayjangde/olxapi/internal/handlers"
	"github.com/abhinayjangde/olxapi/internal/logging"
	middleware "github.com/abhinayjangde/olxapi/internal/middlewares"
	"github.com/rs/cors"
)

func main() {
	// loading env config
	cfg := config.MustLoad()

	// logger setup
	logger, closer, err := logging.New(cfg.LogFile)
	if err != nil {
		slog.Error("logger initialization failed", "err", err)
		os.Exit(1)
	}
	defer closer.Close()
	slog.SetDefault(logger)

	redis, err := db.NewRedisClient(cfg.RedisUrl)
	if err != nil {
		logger.Error("redis initialization failed", "err", err)
		os.Exit(1)
	}
	defer redis.Close()

	db, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		logger.Error("database connection failed", "err", err)
		os.Exit(1)
	}
	mux := http.NewServeMux()
	// wrappedMux := middleware.CorsMiddleware(mux) // for fixing cors policy
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:5500"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            false, // Set to true to print debugging info to stderr
	})

	logger.Info("redis ready")
	logger.Info("postgres database connected")
	logger.Info("starting server", "port", cfg.Port)

	lh := handlers.NewListingHandler(db, redis, logger) // listing handler

	wrappedMux := middleware.RequestId(c.Handler(mux))

	mux.HandleFunc("GET /", handlers.Home)
	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.HandleFunc("GET /listings", lh.List)
	mux.HandleFunc("POST /listings", lh.Create)
	mux.HandleFunc("DELETE /listings/{id}", lh.Delete)

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      wrappedMux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("server stopped", "err", err)
		os.Exit(1)
	}
}
