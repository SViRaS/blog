package handlers

import (
	"blog/internal/config"
	"blog/internal/models"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

type Handler struct {
	DB        *gorm.DB
	Templates *template.Template
	Store     *sessions.CookieStore
}

func New(db *gorm.DB) *Handler {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	return &Handler{
		DB:        db,
		Templates: tmpl,
		Store:     config.Store,
	}
}

func (h *Handler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	h.Templates.ExecuteTemplate(w, "register.html", nil)
}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	h.Templates.ExecuteTemplate(w, "login.html", nil)
}

func (h *Handler) RegisterSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if username == "" || email == "" || password == "" {
		http.Error(w, "Все поля обязательны", http.StatusBadRequest)
		return
	}

	user := models.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		fmt.Println("Ошибка регистрации", err)
		http.Error(w, "Пользователь с такой почтой уже существует", http.StatusConflict)
		return
	}

	fmt.Printf("Пользователь зарегистрирован: %s (ID: %d)\n", username, user.ID)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *Handler) LoginSubmit(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	var user models.User
	if err := h.DB.Where("email = ?", email).First(&user).Error; err != nil {
		http.Error(w, "Неверная почта или пароль", http.StatusUnauthorized)
		return
	}

	if !user.CheckPassword(password) {
		http.Error(w, "Неверная почта или пароль", http.StatusUnauthorized)
		return
	}

	session, _ := h.Store.Get(r, "user-session")
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	session, _ := h.Store.Get(r, "user-session")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
