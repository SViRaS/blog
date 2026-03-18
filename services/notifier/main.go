package main

import (
	"blog/services/notifier/models"
	"encoding/json"
	"log"
	"net/http"
)

func main() {

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

		log.Printf("Excellent, vse rabotaet")
		w.WriteHeader(http.StatusNoContent)
	})

	http.ListenAndServe(":7071", mux)
}
