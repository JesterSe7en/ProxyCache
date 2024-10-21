package main

import (
	"sync"
	"time"
)

type TokenBucket struct {
	maxTokens         int           // total maximum number of tokens
	tokens            int           // current number of tokens
	tokensPerInterval int           // how many tokens to add per refill interval
	refillInterval    time.Duration // how often to refill
	mu                sync.Mutex    // to avoid race conditions
	stopCh            chan struct{} // channel to signal to stop the refilling process
}

func NewTokenBucket(maxTokens, tokensPerInterval int, refillInterval time.Duration) *TokenBucket {
	tb := &TokenBucket{
		maxTokens:         maxTokens,
		tokens:            maxTokens,
		tokensPerInterval: tokensPerInterval,
		refillInterval:    refillInterval,
		stopCh:            make(chan struct{}),
	}
	
	return tb
}

func (tb *TokenBucket) refillTokens() {
	ticker := time.NewTicker(tb.refillInterval)
	// inifinite loop and listen for the ticker to signal "<-ticker.C" ever tb.refillInterval
	// time.  Once recieved, it will start to add tokens to the bucket.
	// when the stopCh is closed, it will stop the refilling process.
	for {
		select {
		case <-ticker.C:
			tb.mu.Lock()
			tb.tokens += tb.tokensPerInterval
			if tb.tokens > tb.maxTokens {
				tb.tokens = tb.maxTokens
			}
			tb.mu.Unlock()
		case <-tb.stopCh:
			return
		}
	}
}

func (tb *TokenBucket) stop() {
	close(tb.stopCh)
}

// TakeToken takes a token from the bucket if available. If the bucket is empty it will return false.
func (tb *TokenBucket) TakeToken() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}
