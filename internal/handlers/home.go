package handlers

import (
	"blog/internal/models"
	"fmt"
	"log"
	"net/http"
)

type PageData struct {
	Title       string
	Heading     string
	Content     string
	Posts       []string
	CurrentUser *models.User
}

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(r, h.DB, h.Store)

	var posts []models.Post
	h.DB.
		Order("created_at DESC").
		Preload("User").
		Preload("Comments.User").
		Find(&posts)

	data := map[string]interface{}{
		"Title":       "Блоги",
		"Posts":       posts,
		"CurrentUser": user,
		"Heading":     "🏠 Главная",
	}

	if user != nil {
		data["Heading"] = fmt.Sprintf("👋 Привет, %s!", user.Username)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := h.Templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Printf("❌ Ошибка шаблона index.html: %v", err)
		http.Error(w, "Ошибка шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
