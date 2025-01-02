package RateLimiters

import (
	"TODO-LIST/Deserializers"
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

// RateLimiter defines a rate limiter with limits per client
type RateLimiter struct {
	clients map[string]*rate.Limiter
	mu      sync.Mutex
	rate    rate.Limit
	burst   int
}

// NewRateLimiter initializes a new RateLimiter
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*rate.Limiter),
		rate:    r,
		burst:   b,
	}
}

// GetLimiter retrieves or creates a rate limiter for a specific client if not already mapped
func (rl *RateLimiter) GetLimiter(clientID string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.clients[clientID]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.clients[clientID] = limiter
	}
	return limiter
}

// Middleware to enforce rate limiting
func (rl *RateLimiter) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		creds, r := Deserializers.ClientAuthenticate(w, r)
		if creds.ClientID == "" {
			http.Error(w, "Client ID required", http.StatusBadRequest)
			return
		}

		limiter := rl.GetLimiter(creds.ClientID)
		fmt.Println("Limiter fetched", limiter)
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		clientID := r.Header.Get("client_id")
		clientSecret := r.Header.Get("client_secret")
		fmt.Println(clientID, clientSecret)
		next.ServeHTTP(w, r)
	})
}
