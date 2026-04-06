package main

import (
	"blog/services/notifier/models"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../../.env")

	telegramURL := os.Getenv("TELEGRAM_NOTIFY_URL")
	if telegramURL == "" {
		telegramURL = "http://localhost:7072/notify"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", 405)
			return
		}

		payload := models.Payload{}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid json", 400)
			return
		}

		if payload.Event == "" {
			http.Error(w, "event is required", 400)
			return
		}

		body, err := json.Marshal(payload)
		if err != nil {
			log.Printf("notifier: marshal payload error: %v", err)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		client := &http.Client{Timeout: 2 * time.Second}
		req, err := http.NewRequest(http.MethodPost, telegramURL, bytes.NewReader(body))
		if err != nil {
			log.Printf("notifier: build request error: %v", err)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("notifier: telegram request error: %v", err)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		_ = resp.Body.Close()
		if resp.StatusCode >= 300 {
			log.Printf("notifier: telegram non-2xx status: %s", resp.Status)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	http.ListenAndServe(":7071", mux)
}
