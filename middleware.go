package main

import (
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type rateLimiter struct {
	sync.Mutex
	requests map[string][]time.Time
}

func newRateLimiter() *rateLimiter {
	rl := &rateLimiter{
		requests: make(map[string][]time.Time),
	}
	// Cleanup routine every 5 minutes
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			rl.Lock()
			now := time.Now()
			for ip, times := range rl.requests {
				var valid []time.Time
				for _, t := range times {
					if now.Sub(t) < time.Minute {
						valid = append(valid, t)
					}
				}
				if len(valid) == 0 {
					delete(rl.requests, ip)
				} else {
					rl.requests[ip] = valid
				}
			}
			rl.Unlock()
		}
	}()
	return rl
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.Lock()
	defer rl.Unlock()

	now := time.Now()
	cutoff := now.Add(-1 * time.Minute)

	times := rl.requests[ip]
	var recent []time.Time
	for _, t := range times {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}

	if len(recent) >= 3 {
		rl.requests[ip] = recent
		return false
	}

	rl.requests[ip] = append(recent, now)
	return true
}

func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func rateLimitMiddleware(rl *rateLimiter, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		if !rl.allow(ip) {
			w.Header().Set("HX-Retarget", "#form-error")
			w.Header().Set("HX-Reswap", "innerHTML")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`<div class="bg-red-500/20 border border-red-500/50 text-red-200 p-4 rounded-xl text-center text-sm font-semibold">
				Bạn đã gửi quá nhiều yêu cầu. Vui lòng thử lại sau 1 phút.
			</div>`))
			return
		}
		next(w, r)
	}
}

func panicRecoveryMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Printf("PANIC: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func loggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			logger.Printf(
				"%s %s %s %d %v",
				r.RemoteAddr,
				r.Method,
				r.URL.Path,
				wrapped.statusCode,
				duration,
			)
		})
	}
}
