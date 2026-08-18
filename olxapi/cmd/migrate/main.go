package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up | down>")
		return
	}

	switch os.Args[1] {
	case "up":
		log.Println("Running migrations up...")
		// Add your migration logic here
	case "down":
		log.Println("Running migrations down...")
		// Add your rollback logic here
	default:
		log.Fatal("usage: migrate <up | down>")
	}
}
