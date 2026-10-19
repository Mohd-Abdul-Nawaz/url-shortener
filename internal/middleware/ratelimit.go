package middleware

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu      sync.Mutex
	clients map[string][]time.Time
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string][]time.Time),
		limit:   limit,
		window:  window,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Clean old entries
	var recent []time.Time
	for _, t := range rl.clients[ip] {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}

	if len(recent) >= rl.limit {
		return false
	}

	recent = append(recent, now)
	rl.clients[ip] = recent
	return true
}

func RateLimit(next http.Handler) http.Handler {
	rl := NewRateLimiter(100, time.Minute) // 100 requests per minute per IP

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.Allow(r.RemoteAddr) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}