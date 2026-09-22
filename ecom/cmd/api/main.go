package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abhinayjangde/ecom/internal/config"
	"github.com/abhinayjangde/ecom/internal/db"
	"github.com/abhinayjangde/ecom/internal/handlers"
	middleware "github.com/abhinayjangde/ecom/internal/middlewares"
	"github.com/abhinayjangde/ecom/internal/utils"
	"github.com/rs/cors"
)

func main() {
	// loading env config
	cfg := config.MustLoad()

	// logger setup
	logger, closer, err := utils.NewLogger(cfg.LogFile)
	if err != nil {
		slog.Error("logger initialization failed", "err", err)
		os.Exit(1)
	}
	defer closer.Close()
	slog.SetDefault(logger)

	redis, err := db.NewRedisClient(cfg.RedisURL)
	if err != nil {
		logger.Error("redis initialization failed", "err", err)
		os.Exit(1)
	}
	defer redis.Close()

	db, err := db.Connect(cfg.DatabaseURL)
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

	lh := handlers.NewListingHandler(db, redis, logger) // listing handler
	uh := handlers.NewUserHandler(db, logger, cfg)      // user handler

	wrappedMux := middleware.RequestId(c.Handler(mux))

	mux.HandleFunc("GET /", handlers.Home)
	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.HandleFunc("GET /listings", lh.List)
	mux.HandleFunc("POST /listings", lh.Create)
	mux.HandleFunc("DELETE /listings/{id}", lh.Delete)

	mux.HandleFunc("POST /auth/register", uh.Create)
	mux.HandleFunc("POST /auth/login", uh.Login)

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      wrappedMux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	// this keeps main thread alive and allows the server to run in the background
	go func() {
		logger.Info("starting server", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server forced to close: ", err.Error(), "\n")
		}
	}()

	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	sig := <-shutdown
	logger.Error("Received signal: ", sig.String(), ". Initiating graceful shutdown...\n")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Graceful shutdown failed:", err.Error(), "\n")
	}

	logger.Info("Server exited cleanly.")
}
