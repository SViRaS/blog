package handlers

import (
	"blog/internal/models"
	"net/http"

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
