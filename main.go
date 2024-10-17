package main

import (
	"flag"
	"os"
)

func main() {
	// Parse command-line arguments
	port := flag.Int("port", -1, "Required. The port on which the server will listen for incoming requests")
	redirectURL := flag.String("redirectURL", "", "Required. The URL of the external service to be proxied")
	flag.Parse()

	if *port == -1 || *redirectURL == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Initialize logger
	err := initLogger()
	defer closeLogger()
	if err != nil {
		LogFatal("Cannot initialize logger", err)
	}

	// Load Redis configuration
	redisConfig, err := loadConfig()
	if err != nil {
		LogFatal("Cannot load Redis config", err)
	}
	if redisConfig == nil {
		LogFatal("Redis config is nil", nil)
	}

	// Connect to Redis
	redisCache := &RedisCache{}
	redisClient, err := connectToRedis(redisConfig)
	if err != nil {
		LogFatal("Cannot connect to Redis", err)
	}
	if redisClient == nil {
		LogFatal("Redis client is nil after connection attempt", nil)
	}
	redisCache.redisClient = redisClient

	// Start the server
	startServer(*port, redisCache, *redirectURL)
}
