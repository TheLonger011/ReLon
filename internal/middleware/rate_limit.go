package middleware

import (
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"sync"
)

type IPLimiter struct {
	mtx      sync.Mutex
	limiters map[string]*rate.Limiter
}

func NewIPLimiter() *IPLimiter {
	limiter := &IPLimiter{limiters: make(map[string]*rate.Limiter)}
	return limiter
}

func (il *IPLimiter) GetLimiter(ip string) *rate.Limiter {
	il.mtx.Lock()
	defer il.mtx.Unlock()

	limiter, exists := il.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(5), 10)
		il.limiters[ip] = limiter
	}
	return limiter
}

func RateLimit(limiter *IPLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			limit := limiter.GetLimiter(ip)
			if !limit.Allow() {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
