package main

import (
	"blog/internal/config"
	"blog/internal/database"
	"blog/internal/handlers"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config.InitStore()

	if err := database.ConnectDB(); err != nil {
		log.Fatal("❌ Не удалось подключиться к БД:", err)
	}

	log.Println("✅ Подключение к БД настроено")

	h := handlers.New(database.DB)

	r := chi.NewRouter()

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.Get("/", h.HomeHandler)
	r.Get("/register", h.RegisterPage)
	r.Post("/register", h.RegisterSubmit)
	r.Get("/login", h.LoginPage)
	r.Post("/login", h.LoginSubmit)
	r.Post("/logout", h.Logout)

	log.Println("🚀 Сервер запущен на http://localhost:7070")
	log.Fatal(http.ListenAndServe(":7070", r))
}
