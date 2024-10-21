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
func startServer(server *http.Server, cache Cache, redirectURL string, tb *TokenBucket) error {
	if server == nil || redirectURL == "" || cache == nil || tb == nil {
		return fmt.Errorf("invalid arguments provided; cannot start server")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		success := tb.TakeToken()

		if !success {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		RootHandler(w, r, cache, redirectURL)
	})

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
	if err != nil {
		LogError("failed to get cached response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(cachedResponse) > 0 {
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
	err = cache.setCachedResponse(hashKey, resData, defaultExpiration)
	if err != nil {
		LogError("failed to cache response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(resData)

	LogInfo(fmt.Sprintf("response sent: %s %s", r.Method, r.URL.Path))
}

// forwardRequest forwards the given request to the given redirect URL.
//
// This function is responsible for creating a new request with the given method,
// redirect URL, and body. It then copies over all header values from the
// original request to the new request. Finally, it uses the http.DefaultClient
// to send the new request and returns the response from the server.
//
// If any step of this process fails, an error is returned. If the response from
// the server is nil, an error is also returned.
//
// forwardRequest is used by the RootHandler to forward incoming requests to the
// given redirect URL.
func forwardRequest(req *http.Request, redirectURL string) (*http.Response, error) {
	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}

	if redirectURL == "" {
		return nil, fmt.Errorf("redirect URL is empty")
	}

	newReq, err := http.NewRequest(req.Method, redirectURL, req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	// Copy over header values from the original request to the new request.
	for name, values := range req.Header {
		for _, value := range values {
			newReq.Header.Add(name, value)
		}
	}

	client := http.DefaultClient
	resp, err := client.Do(newReq)
	if err != nil {
		return nil, fmt.Errorf("failed to forward request: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("response is nil after forwarding request")
	}

	return resp, nil
}
