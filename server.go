package main

import (
	"fmt"
	"net/http"
	"time"
)

type Cache interface {
	getCachedReponse(key string) ([]byte, error)
	setCachedReponse(key string, value []byte, expiration time.Duration) error
}

func startServer(port int, cache Cache, redirectURL string) {
	LogInfo(fmt.Sprintf("Starting server on port %d", port))
	LogInfo(fmt.Sprintf("Redirect URL: %s", redirectURL))

}

func handleRequest(w http.ResponseWriter, r http.Request) {
}
