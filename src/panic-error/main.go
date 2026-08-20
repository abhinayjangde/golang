package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func main3() {

	http.HandleFunc("/users", userHandler)

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		panic(err)
	}

}

func userHandler(w http.ResponseWriter, r *http.Request) {

	// Recover from panic and log the error
	defer func() {
		rec := recover()
		if rec != nil {
			slog.Error("panic recovered", "err", rec, "stack", string(debug.Stack()), "path", r.URL.Path)

			http.Error(w, "User not found", http.StatusNotFound)
		}
	}()

	id := r.URL.Query().Get("id")

	u := getUser(id)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(u.Name + "\n"))
}

type User struct {
	ID   string
	Name string
}

var users = map[string]*User{
	"1": {ID: "1", Name: "Alice"},
	"2": {ID: "2", Name: "Bob"},
}

func getUser(id string) *User {
	return users[id]
}
