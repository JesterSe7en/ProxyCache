package main

import (
	"fmt"
	"net/http"
)

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, expiration int) error
}

func startServer(port int, cache Cache, redirectURL string) {
	LogInfo(fmt.Sprintf("Starting server on port %d", port))
	LogInfo(fmt.Sprintf("Redirect URL: %s", redirectURL))
}

func handleRequest(w http.ResponseWriter, r http.Request) {
}
