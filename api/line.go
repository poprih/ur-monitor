package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/poprih/ur-monitor/db"
	"github.com/poprih/ur-monitor/lib/models"
)

func HandleLine(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate request body
	if len(body) == 0 {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}
	log.Printf("Received webhook: %s", body)
	var event models.LineWebhookEvent

	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Error parsing webhook event: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	database, err := db.ConnectDB()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	for _, e := range event.Events {
		switch e.Type {
		case "follow":
			_, err = database.Exec("INSERT INTO users (line_user_id) VALUES ($1) ON CONFLICT (line_user_id) DO NOTHING", e.Source.UserID)
			if err != nil {
				log.Println("Error inserting user:", err)
				http.Error(w, "Failed to save user", http.StatusInternalServerError)
				return
			}
			fmt.Fprint(w, "User saved successfully")

		case "message":
			// var userID int
			// err = database.QueryRow("SELECT id FROM users WHERE line_user_id = $1", event.LineUserID).Scan(&userID)
			// if err != nil {
			// 	http.Error(w, "User not found", http.StatusNotFound)
			// 	return
			// }

			// _, err = database.Exec("INSERT INTO subscriptions (user_id, unit_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", userID, event.UnitID)
			// if err != nil {
			// 	log.Println("Error inserting subscription:", err)
			// 	http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
			// 	return
			// }

			// fmt.Fprint(w, "Subscription saved successfully")

		default:
			http.Error(w, "Invalid event type", http.StatusBadRequest)
		}
	}
}
