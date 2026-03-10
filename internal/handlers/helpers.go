package handlers

import (
	"blog/internal/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

func GetCurrentUser(r *http.Request, db *gorm.DB, store *sessions.CookieStore) *models.User {
	session, _ := store.Get(r, "user-session")
	userID, ok := session.Values["user_id"].(uint)

	if !ok || userID == 0 {
		return nil
	}

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil
	}

	return &user
}

func GetCurrentPost(r *http.Request, db *gorm.DB) (*models.Post, error) {
	postIDStr := chi.URLParam(r, "id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		return nil, err
	}

	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		return nil, err
	}

	return &post, nil
}
