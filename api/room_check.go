package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/poprih/ur-monitor/db"
	"github.com/poprih/ur-monitor/pkg/line"
)

// URResponse represents the response from UR API
type URResponse struct {
	Count int      `json:"count"`
	Room  []string `json:"room"`
}

// CheckRoomsHandler is an HTTP handler that checks for available units
func CheckRoomsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer " + os.Getenv("CHECK_ROOMS_SECRET") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Only allow scheduled requests (from GitHub Actions)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := checkAndNotifyAvailableRooms()
	if err != nil {
		log.Printf("Error checking rooms: %v", err)
		http.Error(w, fmt.Sprintf("Error checking rooms: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Room check completed successfully")
}

// checkAndNotifyAvailableRooms fetches all active subscriptions and checks for available rooms
func checkAndNotifyAvailableRooms() error {
	// Connect to the database
	database, err := db.ConnectDB()
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer database.Close()

	// Get all units that have active subscriptions
	rows, err := database.Query(`
		SELECT DISTINCT u.unit_name, u.unit_code 
		FROM units u
		JOIN subscriptions s ON u.id = s.unit_id
		JOIN users usr ON s.line_user_id = usr.line_user_id
		WHERE s.deleted_at IS NULL
	`)
	if err != nil {
		return fmt.Errorf("failed to query subscribed units: %w", err)
	}
	defer rows.Close()

	// Process each unit
	for rows.Next() {
		var unitName, unitCode string
		if err := rows.Scan(&unitName, &unitCode); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		// Parse unit_code to get required parameters
		parts := strings.Split(unitCode, "_")
		if len(parts) != 2 {
			log.Printf("Invalid unit_code format: %s", unitCode)
			continue
		}

		shisya := parts[0]
		danchi := parts[1][:3]
		shikibetu := parts[1][3:]

		// Check if this unit has available rooms
		response, err := checkAvailableRooms(shisya, danchi, shikibetu)
		if err != nil {
			log.Printf("Error fetching data for unit %s: %v", unitName, err)
			continue
		}

		// If rooms are available, notify subscribed users
		if response.Count > 0 {
			err = notifySubscribedUsers(database, unitName, response)
			if err != nil {
				log.Printf("Error notifying users for unit %s: %v", unitName, err)
			}
		} else {
			log.Printf("No available rooms for unit %s", unitName)
		}

		// Add a small delay to avoid overwhelming the UR API
		time.Sleep(1 * time.Second)
	}

	return nil
}

// checkAvailableRooms fetches available room data from the UR API
func checkAvailableRooms(shisya, danchi, shikibetu string) (*URResponse, error) {
	baseURL := os.Getenv("UR_API_BASE_URL")
	apiPath := os.Getenv("UR_UNIT_ROOM_CHECK_PATH")
	if baseURL == "" {
		return nil, fmt.Errorf("UR_API_BASE_URL is not set")
	}

	url := fmt.Sprintf("%s%s", baseURL, apiPath)
	postData := fmt.Sprintf("shisya=%s&danchi=%s&shikibetu=%s", shisya, danchi, shikibetu)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(postData)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Origin", "https://www.ur-net.go.jp")
	req.Header.Set("Referer", "https://www.ur-net.go.jp/")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data URResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &data, nil
}

// notifySubscribedUsers notifies all users subscribed to a particular unit
func notifySubscribedUsers(db *sql.DB, unitName string, response *URResponse) error {
	// Find all users subscribed to this unit
	rows, err := db.Query(`
		SELECT usr.line_user_id, usr.reply_token, s.room_types
		FROM users usr
		JOIN subscriptions s ON usr.line_user_id = s.line_user_id
		JOIN units u ON s.unit_id = u.id
		WHERE u.unit_name = $1 AND s.deleted_at IS NULL
	`, unitName)
	if err != nil {
		return fmt.Errorf("failed to query subscribed users: %w", err)
	}
	defer rows.Close()

	channelToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	if channelToken == "" {
		return fmt.Errorf("LINE_CHANNEL_ACCESS_TOKEN is not set")
	}
	
	lineClient := line.NewLineClient(channelToken)

	// Send notification to each subscribed user
	for rows.Next() {
		var userID, replyToken string
		var subscribedRoomTypesJSON []byte
		if err := rows.Scan(&userID, &replyToken, &subscribedRoomTypesJSON); err != nil {
			log.Printf("Error scanning user row: %v", err)
			continue
		}

		// Parse subscribed room types
		var subscribedRoomTypes []string
		if len(subscribedRoomTypesJSON) > 0 {
			if err := json.Unmarshal(subscribedRoomTypesJSON, &subscribedRoomTypes); err != nil {
				log.Printf("Error parsing room types JSON: %v", err)
				continue
			}
		}

		// Check if any of the user's subscribed room types match available rooms
		shouldNotify := false
		if len(subscribedRoomTypes) == 0 {
			// If no room types specified, notify for any availability
			shouldNotify = true
		} else {
			// Check if any of the subscribed room types are available
			for _, subscribedRoomType := range subscribedRoomTypes {
				for _, availableRoom := range response.Room {
					if availableRoom == subscribedRoomType {
						shouldNotify = true
						break
					}
				}
				if shouldNotify {
					break
				}
			}
		}

		if !shouldNotify {
			continue
		}

		// Prepare notification message
		var messageBuilder strings.Builder
		messageBuilder.WriteString(fmt.Sprintf("🔔 *UR %s - 空室通知 / Vacancy Notification*\n\n", unitName))
		messageBuilder.WriteString(fmt.Sprintf("空室数 / Available rooms: %d\n", response.Count))
		messageBuilder.WriteString("\n空室タイプ / Available room types:\n")
		for _, roomType := range response.Room {
			messageBuilder.WriteString(fmt.Sprintf("- %s\n", roomType))
		}

		// Get the property URL from the database
		var propertyURL string
		err = db.QueryRow("SELECT url FROM units WHERE unit_name = $1", unitName).Scan(&propertyURL)
		if err != nil {
			log.Printf("Error getting property URL: %v", err)
		} else {
			fullURL := fmt.Sprintf("https://www.ur-net.go.jp%s", propertyURL)
			messageBuilder.WriteString(fmt.Sprintf("\n物件詳細 / Property details: %s\n", fullURL))
		}

		messageBuilder.WriteString("\n⚠️ ご注意 / Important:\n")
		messageBuilder.WriteString("- 空室は先着順です。お早めにご応募ください。\n")
		messageBuilder.WriteString("- この物件の通知は自動的に解除されます。\n\n")
		messageBuilder.WriteString("- Vacancies are filled on a first-come, first-served basis. Please apply as soon as possible.\n")
		messageBuilder.WriteString("- This property notification will be automatically unsubscribed.\n")

		message := messageBuilder.String()

		// Send push notification using LINE API
		err = lineClient.SendPushMessage(userID, message)
		if err != nil {
			log.Printf("Error sending push message to user %s: %v", userID, err)
			continue
		}

		// After successful notification, unsubscribe the user from this unit
		if err := unsubscribeUser(db, userID, unitName); err != nil {
			log.Printf("Error unsubscribing user %s from unit %s: %v", userID, unitName, err)
			continue
		}
	}

	return nil
}

// unsubscribeUser removes the subscription for a specific user and unit
func unsubscribeUser(db *sql.DB, userID string, unitName string) error {
	_, err := db.Exec(`
		UPDATE subscriptions s
		SET deleted_at = NOW()
		FROM units u
		WHERE s.unit_id = u.id
		AND u.unit_name = $1
		AND s.line_user_id = $2
	`, unitName, userID)
	
	if err != nil {
		return fmt.Errorf("failed to unsubscribe user: %w", err)
	}
	return nil
}
