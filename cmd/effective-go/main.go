package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type RateLimiter struct {
	sync.Mutex
	requests      map[string]int
	maxRequests   int
	resetInterval time.Duration
}

func NewRateLimiter(maxRequests int, resetInterval time.Duration) *RateLimiter {
	return &RateLimiter{
		requests:      make(map[string]int),
		maxRequests:   maxRequests,
		resetInterval: resetInterval,
	}
}

func (rl *RateLimiter) LimitRequest(ip string) bool {
	rl.Lock()
	defer rl.Unlock()

	count, ok := rl.requests[ip]
	if !ok {
		rl.requests[ip] = 1
		return true
	}

	if count >= rl.maxRequests {
		return false
	}

	rl.requests[ip]++
	go func() {
		time.Sleep(rl.resetInterval)
		rl.Lock()
		defer rl.Unlock()
		delete(rl.requests, ip)
	}()

	return true
}

func rateLimitMiddleware(rl *RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if !rl.LimitRequest(ip) {
			w.WriteHeader(http.StatusTooManyRequests)
			log.Printf("Too many requests from %s\n", ip)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func greetingHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Hello, %s!", r.RemoteAddr)
}

func main() {
	rateLimiter := NewRateLimiter(3, 1*time.Second)
	mux := http.NewServeMux()
	mux.Handle("/", rateLimitMiddleware(rateLimiter, http.HandlerFunc(greetingHandler)))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("Server listening on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-shutdown

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to gracefully shutdown server: %v", err)
	}

	log.Println("Server stopped")
}
