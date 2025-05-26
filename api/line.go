package api

import (
	"database/sql"
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

// handleUnsubscribe handles the unsubscribe command
func handleUnsubscribe(db *sql.DB, lineClient *line.LineClient, userID, messageText string, replyToken string) error {
	// Extract mansion name from the message (remove the "-" prefix)
	mansionName := strings.TrimSpace(messageText[1:])
	
	// Check if the mansion exists
	var unitID int
	err := db.QueryRow("SELECT id FROM units WHERE unit_name ILIKE $1", mansionName).Scan(&unitID)
	if err != nil {
		if err == sql.ErrNoRows {
			lineClient.SendReplyMessage(replyToken, line.MessageTemplates.InvalidUnitName)
			return fmt.Errorf("unit not found: %s", mansionName)
		}
		lineClient.SendReplyMessage(replyToken, line.MessageTemplates.InvalidUnitName)
		return err
	}
	
	// Cancel subscription (soft delete)
	_, err = db.Exec(`
		UPDATE subscriptions 
		SET deleted_at = NOW() 
		WHERE line_user_id = $1 AND unit_id = $2 AND deleted_at IS NULL`, 
		userID, unitID)
	if err != nil {
		lineClient.SendReplyMessage(replyToken, line.FormatBilingualMessage(line.MessageTemplates.UnsubscribeError, mansionName))
		return err
	}
	
	// Check if there are other subscribers for this mansion
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM subscriptions WHERE unit_id = $1 AND deleted_at IS NULL", unitID).Scan(&count)
	if err == nil && count == 0 {
		// If no other subscribers, update mansion subscription status
		_, err = db.Exec("UPDATE units SET is_subscribed = FALSE WHERE id = $1", unitID)
		if err != nil {
			log.Println("Error updating unit subscription status:", err)
		}
	}
	
	// Send unsubscribe success message
	lineClient.SendReplyMessage(replyToken, line.FormatBilingualMessage(line.MessageTemplates.UnsubscribeSuccess, mansionName))
	return nil
}

// handleSubscribe handles the subscribe command
func handleSubscribe(db *sql.DB, lineClient *line.LineClient, userID string, parts []string, replyToken string) error {
	var unitName string
	var roomTypes []string
	
	if len(parts) == 1 {
		unitName = parts[0]
	} else if len(parts) == 2 {
		unitName = strings.TrimSpace(parts[0])
		roomTypes = strings.Split(strings.TrimSpace(parts[1]), "&")
	} else {
		lineClient.SendReplyMessage(replyToken, line.MessageTemplates.InvalidFormat)
		return fmt.Errorf("invalid message format")
	}

	// Check if user is premium and subscription count
	var isPremium bool
	var subscriptionCount int
	err := db.QueryRow("SELECT is_premium FROM users WHERE line_user_id = $1", userID).Scan(&isPremium)
	if err != nil {
		if err == sql.ErrNoRows {
			// User doesn't exist yet, create them with default values
			_, insertErr := db.Exec("INSERT INTO users (line_user_id, is_premium) VALUES ($1, FALSE) ON CONFLICT (line_user_id) DO NOTHING", userID)
			if insertErr != nil {
				return insertErr
			}
			isPremium = false // Default to non-premium
		} else {
			return err
		}
	}

	if !isPremium {
		err = db.QueryRow("SELECT COUNT(*) FROM subscriptions WHERE line_user_id = $1 AND deleted_at IS NULL", userID).Scan(&subscriptionCount)
		if err != nil {
			// COUNT(*) should never return no rows, but handle it gracefully
			if err == sql.ErrNoRows {
				subscriptionCount = 0
			} else {
				return err
			}
		}

		if subscriptionCount >= 1 {
			lineClient.SendReplyMessage(replyToken, line.MessageTemplates.SubscriptionLimitReached)
			return fmt.Errorf("subscription limit reached")
		}
	}

	// Check if unit exists
	var unitID int
	err = db.QueryRow("SELECT id FROM units WHERE unit_name ILIKE $1", unitName).Scan(&unitID)
	if err != nil {
		if err == sql.ErrNoRows {
			lineClient.SendReplyMessage(replyToken, line.MessageTemplates.InvalidUnitName)
			return fmt.Errorf("unit not found: %s", unitName)
		}
		lineClient.SendReplyMessage(replyToken, line.MessageTemplates.InvalidUnitName)
		return err
	}

	// Convert room types to JSON array
	roomTypesJSON, err := json.Marshal(roomTypes)
	if err != nil {
		return err
	}

	// Insert subscription
	_, err = db.Exec(`
		INSERT INTO subscriptions (line_user_id, unit_id, room_types, deleted_at) 
		VALUES ($1::text, $2, $3, NULL) 
		ON CONFLICT (line_user_id, unit_id) 
		DO UPDATE SET room_types = $3, deleted_at = NULL`, 
		userID, unitID, roomTypesJSON)
	if err != nil {
		lineClient.SendReplyMessage(replyToken, line.FormatBilingualMessage(line.MessageTemplates.SubscriptionError, unitName))
		return err
	}

	// Create confirmation message
	var confirmationMsg string
	if len(roomTypes) > 0 {
		confirmationMsg = fmt.Sprintf("%s\n%s", 
			line.FormatBilingualMessage(line.MessageTemplates.SubscriptionSuccess, unitName),
			line.FormatBilingualMessage(line.MessageTemplates.SpecifiedRoomTypes, strings.Join(roomTypes, "、")))
	} else {
		confirmationMsg = line.FormatBilingualMessage(line.MessageTemplates.SubscriptionSuccess, unitName)
	}

	// Get all active subscriptions for this user
	rows, err := db.Query(`
		SELECT u.unit_name, s.room_types
		FROM subscriptions s
		JOIN units u ON s.unit_id = u.id
		WHERE s.line_user_id = $1 AND s.deleted_at IS NULL
	`, userID)
	if err != nil {
		return err
	}
	defer rows.Close()
	
	var subscriptions []string
	for rows.Next() {
		var subUnitName string
		var subRoomTypesJSON []byte
		if err := rows.Scan(&subUnitName, &subRoomTypesJSON); err != nil {
			continue
		}

		var subRoomTypes []string
		if len(subRoomTypesJSON) > 0 {
			if err := json.Unmarshal(subRoomTypesJSON, &subRoomTypes); err != nil {
				continue
			}
		}

		if len(subRoomTypes) > 0 {
			subscriptions = append(subscriptions, fmt.Sprintf("%s: %s", subUnitName, strings.Join(subRoomTypes, "、")))
		} else {
			subscriptions = append(subscriptions, subUnitName)
		}
	}

	if len(subscriptions) > 0 {
		confirmationMsg += "\n\n" + line.MessageTemplates.CurrentSubscriptions + "\n" + strings.Join(subscriptions, "\n")
	}
	
	lineClient.SendReplyMessage(replyToken, confirmationMsg)
	return nil
}

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
			
			messageText := strings.TrimSpace(e.Message.Text)
			
			// Handle unsubscribe command
			if strings.HasPrefix(messageText, "-") {
				if err := handleUnsubscribe(database, lineClient, userID, messageText, e.ReplyToken); err != nil {
					log.Println("Error handling unsubscribe:", err)
				}
				return
			}
			
			// Handle subscribe command
			parts := strings.Split(messageText, ":")
			if err := handleSubscribe(database, lineClient, userID, parts, e.ReplyToken); err != nil {
				log.Println("Error handling subscribe:", err)
			}

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
						// COUNT(*) should never return no rows, but handle it gracefully
						if err == sql.ErrNoRows {
							count = 0
						} else {
							log.Println("Error counting other subscribers:", err)
							continue
						}
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
