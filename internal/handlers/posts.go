package handlers

import (
	"blog/internal/config"
	"blog/internal/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func (h *Handler) ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")

	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	var post models.Post
	err = h.DB.
		Preload("User").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User").Order("created_at ASC")
		}).
		First(&post, postID).Error

	if err != nil {
		http.Error(w, "Пост не найден", http.StatusNotFound)
		return
	}

	user := GetCurrentUser(r, h.DB, h.Store)
	if user != nil {
		log.Printf("👤 [ViewPostHandler] Текущий пользователь: %s (ID=%d)", user.Username, user.ID)
	} else {
		log.Printf("👤 [ViewPostHandler] Пользователь не авторизован")
	}

	data := map[string]interface{}{
		"Title":       post.Title,
		"Post":        post,
		"CurrentUser": user,
		"Comments":    post.Comments,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	log.Printf("🎨 [ViewPostHandler] Рендеринг шаблона: view-post.html")

	err = h.Templates.ExecuteTemplate(w, "view-post.html", data)
	if err != nil {
		log.Printf("❌ [ViewPostHandler] КРИТИЧЕСКАЯ ОШИБКА рендеринга шаблона: %v", err)
		log.Printf("📋 [ViewPostHandler] Данные шаблона: %+v", data)

		return
	}

}

func (h *Handler) AllPostsHandler(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post
	h.DB.
		Order("created_at DESC").
		Preload("User").
		Find(&posts)

	data := map[string]interface{}{
		"Title": "Посты",
		"Posts": posts,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.Templates.ExecuteTemplate(w, "all-posts.html", data)
}

func (h *Handler) MyPostsHandler(w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(r, h.DB, config.Store)

	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var posts []models.Post
	h.DB.Where("user_id = ?", user.ID).Order("created_at DESC").Find(&posts)

	data := map[string]interface{}{
		"Title": "Мои посты",
		"User":  user,
		"Posts": posts,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.Templates.ExecuteTemplate(w, "my-posts.html", data)
}

func (h *Handler) CreatePostPage(w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(r, h.DB, config.Store)

	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"Title": "Создать пост",
		"User":  user,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.Templates.ExecuteTemplate(w, "create-post.html", data)
}

func (h *Handler) EditPostPage(w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(r, h.DB, config.Store)

	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postIDStr := chi.URLParam(r, "id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	var post models.Post
	if err := h.DB.Where("id = ? AND user_id = ?", postID, user.ID).First(&post).Error; err != nil {
		http.Error(w, "Пост не найден", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Title": "Редактировать пост",
		"User":  user,
		"Post":  post,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.Templates.ExecuteTemplate(w, "edit-post.html", data)
}

func (h *Handler) CreatePostSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	user := GetCurrentUser(r, h.DB, config.Store)
	if user == nil {
		http.Redirect(w, r, "login", http.StatusSeeOther)
		return
	}

	err := r.ParseMultipartForm(20 << 30)
	if err != nil {
		http.Error(w, "Файл слишком большой (макс 20 МБ)", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	if title == "" || content == "" {
		http.Error(w, "Заголовок и контент обязательны", http.StatusBadRequest)
		return
	}

	var imagePath string

	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		if !isValidImage(header.Header.Get("Content-Type")) {
			http.Error(w, "Только изображения (JPG, PNG, GIF)", http.StatusBadRequest)
			return
		}

		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
		imagePath = "static/uploads/" + filename

		dst, err := os.Create(imagePath)
		if err != nil {
			http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
			return
		}

		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Ошибка записи файла", http.StatusInternalServerError)
			return
		}

		log.Printf("Файл загружен: %s", imagePath)
	}

	post := models.Post{
		Title:   title,
		Content: content,
		Image:   imagePath,
		UserID:  user.ID,
	}

	if err := h.DB.Create(&post).Error; err != nil {
		http.Error(w, "Ошибка сохранения: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard/posts", http.StatusSeeOther)
}

func (h *Handler) EditPostSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := GetCurrentUser(r, h.DB, config.Store)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postIDStr := chi.URLParam(r, "id")
	postID, _ := strconv.ParseUint(postIDStr, 10, 32)

	title := r.FormValue("title")
	content := r.FormValue("content")

	var post models.Post
	if err := h.DB.Where("id = ? AND user_id = ?", postID, user.ID).First(&post).Error; err != nil {
		http.Error(w, "Пост не найден", http.StatusNotFound)
		return
	}

	file, header, err := r.FormFile("image")
	if err == nil {
		if post.Image != "" {
			os.Remove(post.Image)
		}

		if !isValidImage(header.Header.Get("Content-Type")) {
			http.Error(w, "Только изображения (JPG, PNG, GIF)", http.StatusBadRequest)
			return
		}

		defer file.Close()

		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
		imagePath := "static/uploads/" + filename

		dst, err := os.Create(imagePath)
		if err != nil {
			http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		io.Copy(dst, file)
		post.Image = imagePath
	}

	post.Title = title
	post.Content = content

	h.DB.Save(&post)

	http.Redirect(w, r, "/dashboard/posts", http.StatusSeeOther)
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := GetCurrentUser(r, h.DB, config.Store)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postIDStr := chi.URLParam(r, "id")
	postID, _ := strconv.ParseUint(postIDStr, 10, 32)

	var post models.Post
	if err := h.DB.Where("id = ? AND user_id = ?", postID, user.ID).First(&post).Error; err != nil {
		http.Error(w, "Пост не найден", http.StatusNotFound)
		return
	}

	if post.Image != "" {
		os.Remove(post.Image)
		log.Printf("🗑️ Файл удалён: %s", post.Image)
	}

	h.DB.Where("id = ? AND user_id = ?", postID, user.ID).Delete(&models.Post{})

	http.Redirect(w, r, "/dashboard/posts", http.StatusSeeOther)
}

func isValidImage(contentType string) bool {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	return allowedTypes[contentType]
}

func (h *Handler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := GetCurrentUser(r, h.DB, config.Store)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postIDStr := chi.URLParam(r, "id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Комментарий не может быть пустым", http.StatusBadRequest)
		return
	}

	comment := models.Comment{
		Content: content,
		UserID:  user.ID,
		PostID:  uint(postID),
	}

	if err := h.DB.Create(&comment).Error; err != nil {
		http.Error(w, "Ошибка сохранения: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
}
