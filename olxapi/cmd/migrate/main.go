package main

import (
	"log"
	"os"

	"github.com/abhinayjangde/olxapi/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up | down>")
		return
	}
	cfg := config.MustLoad()

	m, err := migrate.New("file://migrations", cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("migrate.new: %v", err)
	}

	switch os.Args[1] {
	case "up":
		log.Println("Running migrations up...")
		// Add your migration logic here
		if err := m.Up(); err != nil {
			log.Fatalf("migrate.up: %v", err)
		}
	case "down":
		log.Println("Running migrations down...")
		// Add your rollback logic here
		if err := m.Steps(-1); err != nil {
			log.Fatalf("migrate.down: %v", err)
		}
	default:
		log.Fatal("usage: migrate <up | down>")
	}
}
