package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var defaultExpiration = time.Hour * 24

type RedisCache struct {
	redisClient *redis.Client
}

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

func (rc *RedisCache) getCachedReponse(key string) ([]byte, error) {
	if rc.redisClient == nil {
		return nil, fmt.Errorf("cannot get cached response: no Redis client provided")
	}

	if key == "" {
		return nil, fmt.Errorf("key cannot be empty")
	}

	val, err := rc.redisClient.Get(context.Background(), key).Result()
	// TODO: again does this return nil or redis.nil on success?
	if err != redis.Nil {
		return nil, fmt.Errorf("failed to get cached response: %w", err)
	}

	return []byte(val), nil
}

func (rc *RedisCache) setCachedResponse(key string, value []byte, expiration time.Duration) error {
	if rc.redisClient == nil {
		return fmt.Errorf("cannot set cached response: no Redis client provided")
	}
	if key == "" || value == nil {
		return fmt.Errorf("url and/or response cannot be empty")
	}

	if expiration == 0 {
		expiration = defaultExpiration
	}

	err := rc.redisClient.Set(context.Background(), key, value, expiration)

	// TODO: does this return nil or redis.nil on success?
	return err.Err()
}
