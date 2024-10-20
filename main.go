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

	if *port == -1 || *redirectURL == "" || !isValidPort(*port) {
		flag.Usage()
		os.Exit(1)
	}

	// Initialize logger
	err := initLogger()
	defer closeLogger()
	if err != nil {
		LogFatal("cannot initialize logger", err)
	}

	// Load Redis configuration
	redisConfig, err := loadConfig()
	if err != nil {
		LogFatal("cannot load Redis config", err)
	}
	if redisConfig == nil {
		LogFatal("redis config is nil", nil)
	}

	// Connect to Redis
	redisCache := &RedisCache{}
	redisClient, err := connectToRedis(redisConfig)
	if err != nil {
		LogFatal("cannot connect to Redis", err)
	}
	if redisClient == nil {
		LogFatal("redis client is nil after connection attempt", nil)
	}
	redisCache.redisClient = redisClient

	// ---------------------- Start HTTP server ----------------------
	// https://dev.to/mokiat/proper-http-shutdown-in-go-3fji
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", *port),
	}

	go func() {
		err := startServer(server, redisCache, *redirectURL)
		if !errors.Is(err, http.ErrServerClosed) {
			LogFatal("cannot start web server", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM) // respond to SIGINT(ctrl+c) and SIGTERM (system asks the program to terminate gracefully)
	<-sigChan                                               // block until a signal is received

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second) // 10 seconds wait for graceful shutdown
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		LogFatal("could not shutdown server gracefully: %v", err)
	}
	LogInfo("server gracefully shutdown.")
}
