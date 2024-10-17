package main

import (
	"fmt"
	"net/http"
	"time"
)

type Cache interface {
	getCachedResponse(key string) ([]byte, error)
	setCachedResponse(key string, value []byte, expiration time.Duration) error
}

// Just pass Cache as interfaces are reference types. i.e. this will be passed by reference
func startServer(port int, cache Cache, redirectURL string) {
	LogInfo(fmt.Sprintf("Starting server on port %d", port))
	LogInfo(fmt.Sprintf("Redirect URL: %s", redirectURL))
}

func handleRequest(w http.ResponseWriter, r http.Request) {
}
