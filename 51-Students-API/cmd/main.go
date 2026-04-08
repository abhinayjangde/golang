package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/abhinayjangde/student/internal/config"
)

func main() {

	// load config
	cfg := config.MustLoad()

	// database setup

	// setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server is running..."))
	})

	// setup server

	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	fmt.Printf("Server started at %s", cfg.HTTPServer.Addr)
	err := server.ListenAndServe()

	if err != nil {
		log.Fatalf("Failed to start server")
	}
}
