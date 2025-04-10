package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/poprih/ur-monitor/db"
	"github.com/poprih/ur-monitor/lib/models"
	"github.com/poprih/ur-monitor/pkg/line"
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

	channelToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	if channelToken == "" {
		http.Error(w, "LINE_CHANNEL_ACCESS_TOKEN is not set", http.StatusInternalServerError)
		return
	}
	
	lineClient := line.NewLineClient(channelToken)

	for _, e := range event.Events {
		switch e.Type {
		case "follow":
			// Store both the user ID and their reply token
			_, err = database.Exec("INSERT INTO users (line_user_id, reply_token) VALUES ($1, $2) ON CONFLICT (line_user_id) DO UPDATE SET reply_token = $2", 
				e.Source.UserID, e.ReplyToken)
			if err != nil {
				log.Println("Error inserting/updating user:", err)
				http.Error(w, "Failed to save user", http.StatusInternalServerError)
				return
			}
			fmt.Fprint(w, "User saved successfully")
			lineClient.SendReplyMessage(e.ReplyToken, line.MessageTemplates.WelcomeMessage)

		case "message":
			// Update the reply token for the user with each message
			_, err = database.Exec("UPDATE users SET reply_token = $1 WHERE line_user_id = $2", 
				e.ReplyToken, e.Source.UserID)
			if err != nil {
				log.Println("Error updating reply token:", err)
			}					
			var userID = e.Source.UserID
			var unitName = strings.TrimSpace(e.Message.Text)

			var isPremium bool
			var subscriptionCount int
			err = database.QueryRow("SELECT is_premium FROM users WHERE line_user_id = $1", userID).Scan(&isPremium)
			if err != nil {
				log.Println("Error checking premium status:", err)
				http.Error(w, "Failed to check user status", http.StatusInternalServerError)
				return
			}

			if !isPremium {
				err = database.QueryRow("SELECT COUNT(*) FROM subscriptions WHERE line_user_id = $1 AND deleted_at IS NULL", userID).Scan(&subscriptionCount)
				if err != nil {
					log.Println("Error counting subscriptions:", err)
					http.Error(w, "Failed to check subscription count", http.StatusInternalServerError)
					return
				}

				if subscriptionCount >= 1 {
					lineClient.SendReplyMessage(e.ReplyToken, line.MessageTemplates.SubscriptionLimitReached)
					return
				}
			}

			// Check if urID exists in the units table
			var unitID int
			err = database.QueryRow("SELECT id FROM units WHERE unit_name ILIKE $1", unitName).Scan(&unitID)
			if err != nil {
				log.Println("Error querying unit:", err)
				lineClient.SendReplyMessage(e.ReplyToken, line.MessageTemplates.InvalidUnitName)
				http.Error(w, "Failed to query unit", http.StatusInternalServerError)
				return
			}

			// Establish a relationship between userID and unitID in the subscriptions table
			_, err = database.Exec("INSERT INTO subscriptions (line_user_id, unit_id, deleted_at) VALUES ($1::text, $2, NULL) ON CONFLICT (line_user_id, unit_id) DO UPDATE SET deleted_at = NULL", userID, unitID)
			if err != nil {
				log.Println("Error inserting subscription:", err)
				lineClient.SendReplyMessage(e.ReplyToken, line.FormatBilingualMessage(line.MessageTemplates.SubscriptionError, unitName))
				http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
				return
			}
			fmt.Fprint(w, "Subscription saved successfully")
			lineClient.SendReplyMessage(e.ReplyToken, line.FormatBilingualMessage(line.MessageTemplates.SubscriptionSuccess, unitName))
		case "unfollow":
			// Check if the user has any subscriptions
			rows, err := database.Query("SELECT unit_id FROM subscriptions WHERE line_user_id = $1", e.Source.UserID)
			if err != nil {
				log.Println("Error querying user subscriptions:", err)
			} else {
				defer rows.Close()
				
				// Process each unit the user is subscribed to
				for rows.Next() {
					var unitID int
					if err := rows.Scan(&unitID); err != nil {
						log.Println("Error scanning unit ID:", err)
						continue
					}
					
					// Check if this unit has other subscribers
					var count int
					err = database.QueryRow("SELECT COUNT(*) FROM subscriptions WHERE unit_id = $1 AND line_user_id != $2", 
						unitID, e.Source.UserID).Scan(&count)
					if err != nil {
						log.Println("Error counting other subscribers:", err)
						continue
					}
					
					// If no other subscribers, set the unit subscription status to false
					if count == 0 {
						_, err = database.Exec("UPDATE units SET is_subscribed = FALSE WHERE id = $1", unitID)
						if err != nil {
							log.Println("Error updating unit subscription status:", err)
						} else {
							log.Printf("Updated unit ID %d as it has no more subscribers", unitID)
						}
					}
				}
				
					// Delete all subscriptions for this user 
				_, err = database.Exec("DELETE FROM subscriptions WHERE line_user_id = $1", e.Source.UserID)
				if err != nil {
					log.Println("Error deleting user subscriptions:", err)
				}
			}
				// Delete the user from the users table
			_, err = database.Exec("DELETE FROM users WHERE line_user_id = $1", e.Source.UserID)
		
			if err != nil {
				log.Println("Error deleting user:", err)
			}			
			
		default:
			http.Error(w, "Invalid event type", http.StatusBadRequest)
		}
	}
}
