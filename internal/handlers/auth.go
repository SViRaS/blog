package handlers

import (
	"blog/internal/auth"
	"blog/internal/config"
	"blog/internal/models"
	"encoding/json"
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

func (h *Handler) APILogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "метод не разрешён"})
		return
	}

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "неверный JSON"})
		return
	}

	if body.Email == "" || body.Password == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "email и password обязательны"})
		return
	}

	var user models.User
	if err := h.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "неверная почта или пароль"})
		return
	}

	if !user.CheckPassword(body.Password) {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "неверная почта или пароль"})
		return
	}

	token, ttl, err := auth.IssueAccessToken(user.ID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "ошибка выдачи токена"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"access_token": token,
		"expires_in":   int(ttl.Seconds()),
	})
}

func respondJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

func (h *Handler) APIRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "метод не разрешён"})
		return
	}

	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "неверный JSON"})
		return
	}

	if body.RefreshToken == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "refresh_token обязателен"})
		return
	}

	userID, err := auth.ValidateRefreshToken(body.RefreshToken)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "неверный refresh_token"})
		return
	}

	token, accessTTL, err := auth.IssueAccessToken(userID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "ошибка выдачи токена"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"access_token":       token,
		"access_expires_in":  int(accessTTL.Seconds()),
		"refresh_token":      body.RefreshToken,
		"refresh_expires_in": int(auth.RefreshTTL().Seconds()),
	})
}
