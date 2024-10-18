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
func startServer(server *http.Server, cache Cache, redirectURL string) error {
	if server == nil || redirectURL == "" || cache == nil {
		LogFatal("invalid arguments provided; cannot start server", nil)
	}

	LogInfo(fmt.Sprintf("starting server on port %s", server.Addr))
	LogInfo(fmt.Sprintf("redirect URL: %s", redirectURL))

	return server.ListenAndServe()
}

func isValidPort(port int) bool {
	return port > 0 && port <= 65535
}
