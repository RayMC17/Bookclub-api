package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func (a *applicationDependencies) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// defer will be called when the stack unwinds
		defer func() {
			// recover() checks for panics
			err := recover()
			if err != nil {
				w.Header().Set("Connection", "close")
				a.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// func (a *applicationDependencies)rateLimit(next http.Handler) http.Handler {
//     // Create a rate limiter. This rate limiter can initially
//     // handle 5 requests. But after that it can only handle 2 per second
//     // (1 request every half-second).  If after handling 5 initial requests
//     // our server, a quarter-second later, receives another request, that
//     // request will be blocked/dropped/queued since we need half-second to be
//     // able to process one new request. The bucket is initially full(5) but
//     // emptied after 5 requests. We need time to refill it.
//     limiter := rate.NewLimiter(2, 5)
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         // Check if this incoming request will be allowed
//         // We will implement rateLimitExceededResponse(w, r) later
//         // the Allow() method tries to remove a token from the bucket
//         // it returns false if the bucket is empty
//         if !limiter.Allow() {
//             a.rateLimitExceededResponse(w, r)
//             return
//         }
//             next.ServeHTTP(w,r)
//         })
//      }

func (a *applicationDependencies) rateLimit(next http.Handler) http.Handler {
	// Define a rate limiter struct
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time // remove map entries that are stale
	}
	var mu sync.Mutex                      // use to synchronize the map
	var clients = make(map[string]*client) // the actual map
	// A goroutine to remove stale entries from the map
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock() // begin cleanup
			// delete any entry not seen in three minutes
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock() // finish clean up
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.config.limiter.enabled {
			// get the IP address
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				a.serverErrorResponse(w, r, err)
				return
			}

			mu.Lock() // exclusive access to the map
			// check if ip address already in map, if not add it
			_, found := clients[ip]
			if !found {
				clients[ip] = &client{limiter: rate.NewLimiter(
					rate.Limit(a.config.limiter.rps),
					a.config.limiter.burst),
				}
			}
			// Update the last seem for the client
			clients[ip].lastSeen = time.Now()

			// Check the rate limit status
			if !clients[ip].limiter.Allow() {
				mu.Unlock() // no longer need exclusive access to the map
				a.rateLimitExceededResponse(w, r)
				return
			}

			mu.Unlock()
		} // others are free to get exclusive access to the map
		next.ServeHTTP(w, r)
	})

}
