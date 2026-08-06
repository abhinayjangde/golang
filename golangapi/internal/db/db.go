package db

import (
	"context"
	"fmt"
	"time"

	"github.com/abhinayjangde/todo/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(cfg config.Config) (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.DatabaseURL)

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return nil, nil, fmt.Errorf("Database connection failed")
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, fmt.Errorf("Database ping failed")
	}
	database := client.Database("mydb")

	return client, database, err
}

func Disconnect(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return client.Disconnect(ctx)
}
