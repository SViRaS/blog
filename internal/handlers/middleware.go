package handlers

import (
	"blog/internal/config"
	"net/http"
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
