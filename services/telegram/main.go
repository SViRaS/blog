package main

import (
	"blog/services/telegram/internal/adapters"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type notifyPayload struct {
	Event          string    `json:"event"`
	PostID         uint      `json:"post_id"`
	CommentID      uint      `json:"comment_id"`
	AuthorID       uint      `json:"author_id"`
	AuthorUsername string    `json:"author_username"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

func main() {
	_ = godotenv.Load("../../.env")

	botToken := os.Getenv("TG_BOT_TOKEN")
	chatIDStr := os.Getenv("TG_CHAT_ID")
	if botToken == "" {
		log.Fatal("TG_BOT_TOKEN is required")
	}
	if chatIDStr == "" {
		log.Fatal("TG_CHAT_ID is required (e.g. TG_CHAT_ID=123456789)")
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Fatalf("invalid TG_CHAT_ID: %v", err)
	}

	telegramAdapter, err := adapters.NewTelegramAdapter(botToken, chatID)
	if err != nil {
		log.Fatalf("failed to init telegram adapter: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		payload := notifyPayload{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if payload.Event == "" {
			http.Error(w, "event is required", http.StatusBadRequest)
			return
		}

		var message string
		switch payload.Event {
		case "post_created":
			username := payload.AuthorUsername
			if username == "" {
				username = fmt.Sprintf("user#%d", payload.AuthorID)
			}
			if payload.Content != "" {
				message = fmt.Sprintf("📝 Новый пост от @%s: %s", username, payload.Content)
			} else {
				message = fmt.Sprintf("📝 Новый пост от @%s (id=%d)", username, payload.PostID)
			}
		case "comment_created":
			username := payload.AuthorUsername
			if username == "" {
				username = fmt.Sprintf("user#%d", payload.AuthorID)
			}
			if payload.Content != "" {
				message = fmt.Sprintf("💬 Новый комментарий от @%s: %s", username, payload.Content)
			} else {
				message = fmt.Sprintf("💬 Новый комментарий от @%s (comment_id=%d)", username, payload.CommentID)
			}
		default:
			message = fmt.Sprintf("🔔 Событие: %s", payload.Event)
		}

		if err := telegramAdapter.SendMessage(r.Context(), message); err != nil {
			log.Printf("telegram: send error: %v", err)
			http.Error(w, "telegram send failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	log.Printf("telegram service listening on %s", ":7072")
	if err := http.ListenAndServe(":7072", mux); err != nil {
		log.Fatalf("telegram service error: %v", err)
	}
}
