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

// ResponseItem represents a property unit from the UR API response
type ResponseItem struct {
	Name           string `json:"name"`
	Floortype      string `json:"floor"`
	RentNormal     string `json:"rent_normal"`
	RoomDetailLink string `json:"roomDetailLink"`
}

// CheckRoomsHandler is an HTTP handler that checks for available units
func CheckRoomsHandler(w http.ResponseWriter, r *http.Request) {
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
		SELECT DISTINCT u.unit_name 
		FROM units u
		JOIN subscriptions s ON u.id = s.unit_id
		JOIN users usr ON s.line_user_id = usr.line_user_id
		WHERE s.deleted_at IS NULL AND usr.active = true
	`)
	if err != nil {
		return fmt.Errorf("failed to query subscribed units: %w", err)
	}
	defer rows.Close()

	// Process each unit (danchi)
	for rows.Next() {
		var danchiID string
		if err := rows.Scan(&danchiID); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		// Check if this danchi has available rooms
		availableRooms, err := checkAvailableRooms(danchiID)
		if err != nil {
			log.Printf("Error fetching data for danchi %s: %v", danchiID, err)
			continue
		}

		// If rooms are available, notify subscribed users
		if len(availableRooms) > 0 {
			err = notifySubscribedUsers(database, danchiID, availableRooms)
			if err != nil {
				log.Printf("Error notifying users for danchi %s: %v", danchiID, err)
			}
		} else {
			log.Printf("No available rooms for danchi %s", danchiID)
		}

		// Add a small delay to avoid overwhelming the UR API
		time.Sleep(2 * time.Second)
	}

	return nil
}

// checkAvailableRooms fetches available room data from the UR API
func checkAvailableRooms(danchi string) ([]ResponseItem, error) {
	url := "https://chintai.r6.ur-net.go.jp/chintai/api/bukken/detail/detail_bukken_room/"
	postData := fmt.Sprintf("rent_low=&rent_high=&floorspace_low=&floorspace_high=&shisya=80&danchi=%s&shikibetu=0&newBukkenRoom=&orderByField=0&orderBySort=0&pageIndex=0&pageIndex=0&sp=", danchi)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(postData)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []ResponseItem
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	// Complete the room detail links with the domain
	for i := range data {
		data[i].RoomDetailLink = "https://www.ur-net.go.jp" + data[i].RoomDetailLink
	}

	return data, nil
}

// notifySubscribedUsers notifies all users subscribed to a particular danchi
func notifySubscribedUsers(db *sql.DB, danchiID string, availableRooms []ResponseItem) error {
	// Find all users subscribed to this danchi
	rows, err := db.Query(`
		SELECT usr.line_user_id, usr.reply_token
		FROM users usr
		JOIN subscriptions s ON usr.line_user_id = s.line_user_id
		JOIN units u ON s.unit_id = u.id
		WHERE u.unit_name = $1 AND s.deleted_at IS NULL AND usr.active = true
	`, danchiID)
	if err != nil {
		return fmt.Errorf("failed to query subscribed users: %w", err)
	}
	defer rows.Close()

	// Prepare notification message
	var messageBuilder strings.Builder
	messageBuilder.WriteString(fmt.Sprintf("🔔 *Notification for Danchi %s*\n\n", danchiID))
	messageBuilder.WriteString(fmt.Sprintf("Found %d available units:\n\n", len(availableRooms)))

	for _, room := range availableRooms {
		messageBuilder.WriteString(fmt.Sprintf("🏠 Name: %s\n🏢 Floor: %s\n💰 Rent: %s\n🔗 Link: %s\n\n",
			room.Name, room.Floortype, room.RentNormal, room.RoomDetailLink))
	}

	message := messageBuilder.String()

	channelToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	if channelToken == "" {
		return fmt.Errorf("LINE_CHANNEL_ACCESS_TOKEN is not set")
	}
		
	lineClient := line.NewClient(channelToken)

	// Send notification to each subscribed user
	for rows.Next() {
		var userID, replyToken string
		if err := rows.Scan(&userID, &replyToken); err != nil {
			log.Printf("Error scanning user row: %v", err)
			continue
		}

		// Send push notification using LINE API
		err = lineClient.SendPushMessage(userID, message)
		if err != nil {
			log.Printf("Error sending push message to user %s: %v", userID, err)
			continue
		}

		// After successful notification, unsubscribe the user from this danchi
		if err := unsubscribeUser(db, userID, danchiID); err != nil {
			log.Printf("Error unsubscribing user %s from danchi %s: %v", userID, danchiID, err)
			continue
		}
	}

	// After all users are notified and unsubscribed, delete the danchi
	if err := deleteDanchi(db, danchiID); err != nil {
		log.Printf("Error deleting danchi %s: %v", danchiID, err)
	}

	return nil
}

// unsubscribeUser removes the subscription for a specific user and danchi
func unsubscribeUser(db *sql.DB, userID string, danchiID string) error {
	_, err := db.Exec(`
		UPDATE subscriptions s
		SET deleted_at = NOW()
		FROM units u
		WHERE s.unit_id = u.id
		AND u.unit_name = $1
		AND s.line_user_id = $2
	`, danchiID, userID)
	
	if err != nil {
		return fmt.Errorf("failed to unsubscribe user: %w", err)
	}
	return nil
}

// deleteDanchi removes the danchi from the units table
func deleteDanchi(db *sql.DB, danchiID string) error {
	_, err := db.Exec(`
		DELETE FROM units
		WHERE unit_name = $1
	`, danchiID)
	
	if err != nil {
		return fmt.Errorf("failed to delete danchi: %w", err)
	}
	return nil
}
