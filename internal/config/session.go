package config

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

func InitStore() {
	key := os.Getenv("SESSION_KEY")

	if key == "" {
		log.Println("SESSION_KEY не задан, используется dev-ключ")
		key = "dev-key-must-be-at-least-32-characters-long-123456"
	}

	Store = sessions.NewCookieStore([]byte(key))

	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}
