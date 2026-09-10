package db

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(redistUrl string) (*redis.Client, error) {

	opt, err := redis.ParseURL(redistUrl)
	if err != nil {
		return nil, fmt.Errorf("redis: %w", err)
	}
	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = client.Ping(ctx).Result()

	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return client, nil
}
