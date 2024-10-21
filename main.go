package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Parse command-line arguments
	port := flag.Int("port", -1, "Required. The port on which the server will listen for incoming requests")
	redirectURL := flag.String("redirectURL", "", "Required. The URL of the external service to be proxied")
	flag.Parse()

	// TODO: santatize redirectURL

	if *port == -1 || *redirectURL == "" || !isValidPort(*port) {
		flag.Usage()
		os.Exit(1)
	}

	// Initialize logger
	err := initLogger()
	if err != nil {
		LogError("cannot initialize logger", err)
		os.Exit(1)
	}
	defer closeLogger()

	// Load Redis configuration
	redisConfig, err := loadConfig()
	if err != nil {
		LogError("cannot load Redis config", err)
		os.Exit(1)
	}

	if redisConfig == nil {
		LogError("redis config is nil", nil)
		os.Exit(1)
	}

	// Connect to Redis
	redisCache := &RedisCache{nil}
	redisClient, err := connectToRedis(redisConfig)
	if err != nil || redisClient == nil {
		LogError("cannot connect to Redis", err)
		os.Exit(1)
	}
	redisCache.redisClient = redisClient

	// Start HTTP server
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", *port),
	}

	tb := NewTokenBucket(100, 10, 1*time.Second)
	go tb.refillTokens()
	defer tb.stop()

	go func() {
		err := startServer(server, redisCache, *redirectURL, tb)
		if !errors.Is(err, http.ErrServerClosed) {
			LogError("cannot start web server", err)
			return
		}
	}()

	// Handle system signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Shutdown the server gracefully
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		LogError("could not shutdown server gracefully: %v", err)
		os.Exit(1)
	}
	LogInfo("server gracefully shutdown.")
}
