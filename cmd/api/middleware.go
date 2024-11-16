package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Middleware: Panic Recovery
func (a *applicationDependencies) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				a.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Middleware: Rate Limiting
func (a *applicationDependencies) rateLimit(next http.Handler) http.Handler {
	clients := make(map[string]*client)
	var mu sync.Mutex

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, c := range clients {
				if time.Since(c.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.config.limiter.enabled {
			next.ServeHTTP(w, r)
			return
		}

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			a.serverErrorResponse(w, r, err)
			return
		}

		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &client{
				limiter: rate.NewLimiter(rate.Limit(a.config.limiter.rps), a.config.limiter.burst),
			}
		}

		clients[ip].lastSeen = time.Now()
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			a.rateLimitExceededResponse(w, r)
			return
		}
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

// Middleware: Logging
func (a *applicationDependencies) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		a.logger.Info("started request", "method", r.Method, "url", r.URL.String())
		next.ServeHTTP(w, r)
		a.logger.Info("completed request", "method", r.Method, "url", r.URL.String(), "duration", time.Since(start))
	})
}

// Middleware: Authentication (Stub)
func (a *applicationDependencies) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			a.unauthorizedResponse(w, r)
			return
		}
		// Extract token and validate it
		//token := strings.TrimPrefix(authHeader, "Bearer ")
		// Here, add logic to validate token
		// If invalid:
		//authHeader := r.Header.get("Authorization")
		//if authHeather == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		a.unauthorizedResponse(w, r)
		 return
		

		// If token is valid, proceed
		next.ServeHTTP(w, r)
	})
}

// Unauthorized Response Helper
func (a *applicationDependencies) unauthorizedResponse(w http.ResponseWriter, r *http.Request) {
	message := "You are not authorized to access this resource"
	a.errorResponseJSON(w, r, http.StatusUnauthorized, message)
}
