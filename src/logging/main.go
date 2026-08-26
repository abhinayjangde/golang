package main

import (
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Hello endpoint was called")

	w.Write([]byte("Hello, World!"))
}

func main() {
	http.HandleFunc("/hello", helloHandler)

	log.Println("Server starting on port 8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}
