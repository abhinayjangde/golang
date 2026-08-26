package main

import (
	"net/http"

	"github.com/abhinayjangde/shorturl/internal/handler"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handler.Healthz)

	srv := http.Server{
		Addr:    ":9000",
		Handler: mux,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
