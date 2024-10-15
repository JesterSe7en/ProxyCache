package main

import (
	"fmt"
	"os"
)

type Config struct {
	redisURL      string
	redisPassword string
	defaultPort   int
	databaseIndex int
}

func loadConfig() (*Config, error) {
	redisURL := os.Getenv("REDIS_URL")

	if redisURL == "" {
		return nil, fmt.Errorf("REDIS_URL must be set")
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	defaultPort := 6379

	return &Config{
		redisURL, redisPassword, defaultPort, 0,
	}, nil
}
