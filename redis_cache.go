package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
	"strings"
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

	LogInfo(fmt.Sprintf("attempting to connect to Redis at %s:%d.", redisConfig.redisURL, redisConfig.defaultPort))

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	LogInfo("connected to Redis instance.")

	return client, nil // Return the connected Redis client
}

func getHashKey(req *http.Request) string {
	// use the SHA256 algo to generate hash
	// use in combination of the request method, path, and query (sorted)

	// 1. Use the request method
	method := req.Method

	// 2. Use the full URL path
	path := req.URL.Path

	// 3. Sort and concatenate query parameters
	query := req.URL.Query()
	var queryParams []string
	for key, values := range query {
		for _, value := range values {
			queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, value))
		}
	}
	sort.Strings(queryParams)
	queryString := strings.Join(queryParams, "&")

	// Combine method, path, and query string to form the raw cache key
	rawKey := fmt.Sprintf("%s_%s?%s", method, path, queryString)

	// Optional: Hash the key to avoid overly long cache keys
	hash := sha256.New()
	hash.Write([]byte(rawKey))
	hashedKey := hex.EncodeToString(hash.Sum(nil))

	return hashedKey

}

// (rc *RedisCache), this is called a pointer reciver
// the RedisCache refenerece is the original and not a copy
func (rc *RedisCache) getCachedResponse(key string) ([]byte, error) {
	if rc.redisClient == nil {
		return nil, fmt.Errorf("cannot get cached response: no Redis client provided")
	}

	if key == "" {
		return nil, fmt.Errorf("key cannot be empty")
	}

	val, err := rc.redisClient.Get(context.Background(), key).Result()
	if err != nil {
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
