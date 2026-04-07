package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func ConnectDB() {
	// MongoDB connection URI
	clientOptions := options.Client().ApplyURI("mongodb://root:password@localhost:27017")

	// Create context (timeout safety)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	// Ping database (check connection)
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB not reachable:", err)
	}

	DB = client

	log.Println("✅ Connected to MongoDB")
}
