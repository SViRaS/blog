package handlers

import (
	"blog/internal/config"
	"net/http"
	"sync"
	"time"
)

var (
	rateMu  sync.Mutex
	clients = make(map[string]*client)
)

type client struct {
	count       int
	windowStart time.Time
}

const (
	rateLimit  = 5
	rateWindow = 3 * time.Second
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := config.Store.Get(r, "user-session")
		if session.Values["user_id"] == nil {
			http.Redirect(w, r, "/login", 400)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		now := time.Now()

		rateMu.Lock()
		c, ok := clients[ip]
		if !ok {
			c = &client{
				count:       0,
				windowStart: now,
			}
			clients[ip] = c
		}

		if now.Sub(c.windowStart) > rateWindow {
			c.windowStart = now
			c.count = 1
		}

		c.count++

		if c.count > rateLimit {
			rateMu.Unlock()
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte("Rate limit exceeded"))
			return
		}

		rateMu.Unlock()

		next.ServeHTTP(w, r)
	})
}
