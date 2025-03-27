package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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
			var userID = e.Source.UserID
			var urID = e.Message.Text

			// Check if urID exists in the units table
			var unitID int
			err = database.QueryRow("SELECT id FROM units WHERE unit_name = $1", urID).Scan(&unitID)
			if err != nil {
				// If not found, insert a new record into the units table
				err = database.QueryRow("INSERT INTO units (unit_name) VALUES ($1) RETURNING id", urID).Scan(&unitID)
				if err != nil {
					log.Println("Error inserting unit:", err)
					http.Error(w, "Failed to create unit", http.StatusInternalServerError)
					return
				}
			}

			// Establish a relationship between userID and unitID in the subscriptions table
			_, err = database.Exec("INSERT INTO subscriptions (user_id, unit_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", userID, unitID)
			if err != nil {
				log.Println("Error inserting subscription:", err)
				http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
				return
			}

			fmt.Fprint(w, "Subscription saved successfully")
			ReplyToLineMessage(e.ReplyToken, "You have successfully subscribed to UR "+urID)
		default:
			http.Error(w, "Invalid event type", http.StatusBadRequest)
		}
	}
}

// ReplyToLineMessage sends a reply message to a LINE user using the LINE Messaging API.
func ReplyToLineMessage(replyToken, message string) error {
	lineAPIURL := "https://api.line.me/v2/bot/message/reply"
	channelToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	if channelToken == "" {
		return fmt.Errorf("LINE_CHANNEL_ACCESS_TOKEN is not set")
	}

	payload := map[string]interface{}{
		"replyToken": replyToken,
		"messages": []map[string]string{
			{
				"type": "text",
				"text": message,
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", lineAPIURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+channelToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("LINE API error: %s", string(body))
	}

	return nil
}
