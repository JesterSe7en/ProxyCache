package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var defaultExpiration = time.Hour * 24

func connectToRedis(redisConfig *Config) (*redis.Client, error) {
	if redisConfig == nil {
		return nil, fmt.Errorf("no Redis config provided")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.redisURL, redisConfig.defaultPort),
		Password: redisConfig.redisPassword,
		DB:       redisConfig.databaseIndex,
	})

	LogInfo(fmt.Sprintf("Attempting to connect to Redis at %s:%d.", redisConfig.redisURL, redisConfig.defaultPort))

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	LogInfo("Connected to Redis instance.")

	return client, nil // Return the connected Redis client
}

func getCachedReponse(redisClient *redis.Client, url string) ([]byte, error) {
	if redisClient == nil {
		return nil, fmt.Errorf("cannot get cached response: no Redis client provided")
	}

	if url == "" {
		return nil, fmt.Errorf("url cannot be empty")
	}

	val, err := redisClient.Get(context.Background(), url).Result()
	// TODO: again does this return nil or redis.nil on success?
	if err != redis.Nil {
		return nil, fmt.Errorf("failed to get cached response: %w", err)
	}

	return []byte(val), nil
}

func setCachedResponse(redisClient *redis.Client, url string, response []byte, expiration time.Duration) error {
	if redisClient == nil {
		return fmt.Errorf("cannot set cached response: no Redis client provided")
	}
	if url == "" || response == nil {
		return fmt.Errorf("url and/or response cannot be empty")
	}

	if expiration == 0 {
		expiration = defaultExpiration
	}

	err := redisClient.Set(context.Background(), url, response, expiration)

	// TODO: does this return nil or redis.nil on success?
	return err.Err()
}
