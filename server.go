package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Cache interface {
	getCachedResponse(key string) ([]byte, error)
	setCachedResponse(key string, value []byte, expiration time.Duration) error
}

// Just pass Cache as interfaces are reference types. i.e. this will be passed by reference
func startServer(server *http.Server, cache Cache, redirectURL string) error {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		RootHandler(w, r, cache, redirectURL)
	})
	if server == nil || redirectURL == "" || cache == nil {
		return fmt.Errorf("invalid arguments provided; cannot start server")
	}

	LogInfo(fmt.Sprintf("starting server on port %s", server.Addr))
	LogInfo(fmt.Sprintf("redirect URL: %s", redirectURL))

	return server.ListenAndServe()
}

func isValidPort(port int) bool {
	return port > 0 && port <= 65535
}

func RootHandler(w http.ResponseWriter, r *http.Request, cache Cache, redirectURL string) {
	if cache == nil || redirectURL == "" {
		LogError("invalid arguments provided; cannot handle request", nil)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	LogInfo(fmt.Sprintf("request received: %s %s", r.Method, r.URL.Path))

	// get hash key
	hashKey := getHashKey(r)
	// check redis cache
	cachedResponse, err := cache.getCachedResponse(hashKey)
	if err == nil && len(cachedResponse) > 0 {
		LogInfo("cache hit: returning cached response")
		w.Write(cachedResponse)
		return
	}

	LogInfo("cache miss: forwarding request to redirect URL")
	res, err := forwardRequest(r, redirectURL)
	if err != nil {
		LogError("failed to forward request", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		LogError("failed to read response body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	LogInfo("forward request complete; caching response")
	cache.setCachedResponse(hashKey, resData, defaultExpiration)

	w.Write(resData)

	LogInfo(fmt.Sprintf("response sent: %s %s", r.Method, r.URL.Path))
}

func forwardRequest(req *http.Request, redirectURL string) (*http.Response, error) {
	req, err := http.NewRequest(req.Method, redirectURL, req.Body)
	if err != nil {
		return nil, err
	}

	// Copy over header values to new request
	for name, values := range req.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
