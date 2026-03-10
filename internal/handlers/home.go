package handlers

import (
	"blog/internal/models"
	"fmt"
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

	var comments []models.Comment
	h.DB.
		Order("created_at DESC").
		Preload("User").
		Find(&comments)

	var posts []models.Post
	h.DB.
		Order("created_at DESC").
		Preload("User").
		Find(&posts)

	data := map[string]interface{}{
		"Title":    "Посты",
		"Posts":    posts,
		"Comments": comments,
		"Heading":  "🏠 Добро пожаловать!",
	}

	if user != nil {
		data["Heading"] = fmt.Sprintf("👋 Привет, %s!", user.Username)
		data["Content"] = fmt.Sprintf("Вы вошли как %s (%s)", user.Username, user.Email)
		data["CurrentUser"] = user
	} else {
		data["Content"] = "Вы не авторизованы. <a href=\"/login\">Войти</a> или <a href=\"/register\">Зарегистрироваться</a>"
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := h.Templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, "Ошибка шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
