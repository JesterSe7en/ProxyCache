package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func connectToRedis(redisConfig *Config) (*redis.Client, error) {
	if redisConfig == nil {
		return nil, fmt.Errorf("no Redis config provided")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.redisURL,
		Password: redisConfig.redisPassword,
		DB:       redisConfig.databaseIndex,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil // Return the connected Redis client
}

func getCachedReponse(url string) ([]byte, error) {
	return nil, nil
}

func setCachedResponse(url string, response []byte, expiration time.Duration) error {
	return nil
}
