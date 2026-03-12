package main

import (
	"blog/internal/config"
	"blog/internal/database"
	"blog/internal/handlers"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config.InitStore()

	if err := database.ConnectDB(); err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}

	log.Println("Подключение к БД настроено")

	h := handlers.New(database.DB)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.Get("/", h.HomeHandler)
	r.Get("/register", h.RegisterPage)
	r.Post("/register", h.RegisterSubmit)
	r.Get("/login", h.LoginPage)
	r.Post("/login", h.LoginSubmit)
	r.Post("/logout", h.Logout)

	r.Route("/post/{id}", func(r chi.Router) {
		r.Get("/", h.ViewPostHandler)
		r.Post("/comment", h.CreateCommentHandler)

		r.Route("/comment/{commentID}", func(r chi.Router) {
			r.Use(handlers.RequireAuth)
			r.Post("/edit", h.EditCommentHandler)
			r.Post("/delete", h.DeleteCommentHandler)
		})
	})

	r.Route("/dashboard", func(r chi.Router) {
		r.Use(handlers.RequireAuth)

		r.Get("/posts", h.MyPostsHandler)
		r.Get("/create", h.CreatePostPage)
		r.Post("/create", h.CreatePostSubmit)
		r.Get("/post/{id}/edit", h.EditPostPage)
		r.Post("/post/{id}/edit", h.EditPostSubmit)
		r.Post("/post/{id}/delete", h.DeletePost)
	})

	log.Println("🚀 Сервер запущен на http://localhost:7070")
	log.Fatal(http.ListenAndServe(":7070", r))
}