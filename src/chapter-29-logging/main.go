package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := 10

	slog.Info(
		"Fetching user",
		"user_id", userID,
	)

	user := User{
		ID:   userID,
		Name: "Abhi",
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(user)

	slog.Info(
		"User fetched successfully",
		"user_id", user.ID,
	)
}

func main() {
	http.HandleFunc("GET /users", getUserHandler)

	slog.Info("Server started", "port", 8080)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		slog.Error(
			"Server failed to start",
			"error", err,
		)
	}
}
