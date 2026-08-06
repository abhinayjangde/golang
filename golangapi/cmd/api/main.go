package main

import (
	"fmt"
	"log"

	"github.com/abhinayjangde/todo/internal/config"
	"github.com/abhinayjangde/todo/internal/db"
	"github.com/abhinayjangde/todo/internal/server"
)

func main() {

	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("Config error")
	}

	client, database, err := db.Connect(cfg)

	if err != nil {
		log.Fatalf("Db error")
	}

	defer func() {
		if err := db.Disconnect(client); err != nil {
			log.Printf("Db disconnected %v", err)
		}
	}()

	router := server.NewRouter(database)

	addr := fmt.Sprintf(":%s", cfg.Port)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed to start")
	}
}
