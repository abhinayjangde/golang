package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func Redis(redistUrl string) (*redis.Client, error) {

	opt, err := redis.ParseURL(redistUrl)
	if err != nil {
		return nil, fmt.Errorf("redis: %w", err)
	}
	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pong, err := client.Ping(ctx).Result()

	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Println("redis connected:", pong)
	return client, nil
}
