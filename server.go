package main

import (
	"net/http"
)

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, expiration int) error
}

func startServer(port int, cache Cache, redirectURL string) {
}

func handleRequest(w http.ResponseWriter, r http.Request) {
}
