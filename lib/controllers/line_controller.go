package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Event struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
	Source    struct {
		Type    string `json:"type"`
		UserID  string `json:"userId"`
		GroupID string `json:"groupId,omitempty"`
		RoomID  string `json:"roomId,omitempty"`
	} `json:"source"`
	ReplyToken string `json:"replyToken,omitempty"`
	Message    struct {
		Type string `json:"type"`
		ID   string `json:"id"`
		Text string `json:"text,omitempty"`
	} `json:"message,omitempty"`
}

// LineWebhookEvent represents the structure of a LINE webhook event
type LineWebhookEvent struct {
	Destination string `json:"destination"`
	Events      []Event
}

// LineReplyMessage represents the structure for sending a reply to LINE
type LineReplyMessage struct {
	ReplyToken string `json:"replyToken"`
	Messages   []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"messages"`
}

const lineReplyAPI = "https://api.line.me/v2/bot/message/reply"

// HandleLineWebhook processes incoming LINE webhook events
func HandleLineWebhook(w http.ResponseWriter, r *http.Request) {
	// Validate request method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

	// Parse webhook event
	var webhookEvent LineWebhookEvent
	if err := json.Unmarshal(body, &webhookEvent); err != nil {
		log.Printf("Error parsing webhook event: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Process events
	for _, event := range webhookEvent.Events {
		switch event.Type {
		case "message":
			if event.Message.Type == "text" {
				err := processTextMessage(event)
				if err != nil {
					log.Printf("Error processing text message: %v", err)
				}
			}
		case "follow":
			// Handle when a user adds the bot as a friend
			sendReplyMessage(event.ReplyToken, "Thanks for adding me! I'm here to help.")
		case "unfollow":
			// Handle when a user blocks the bot
			log.Printf("User %s unfollowed the bot", event.Source.UserID)
		default:
			log.Printf("Unhandled event type: %s", event.Type)
		}
	}

	w.WriteHeader(http.StatusOK)
}

// processTextMessage handles incoming text messages
func processTextMessage(event Event) error {
	// Example: Echo the received message with some processing
	receivedText := event.Message.Text
	responseText := fmt.Sprintf("You said: %s", receivedText)

	return sendReplyMessage(event.ReplyToken, responseText)
}

// sendReplyMessage sends a reply to the LINE messaging API
func sendReplyMessage(replyToken, message string) error {
	// Retrieve LINE channel access token
	accessToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	if accessToken == "" {
		return fmt.Errorf("LINE_CHANNEL_ACCESS_TOKEN is not set")
	}

	// Prepare reply message
	replyMsg := LineReplyMessage{
		ReplyToken: replyToken,
		Messages: []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		}{
			{
				Type: "text",
				Text: message,
			},
		},
	}

	// Convert to JSON
	payload, err := json.Marshal(replyMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal reply message: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", lineReplyAPI, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send reply: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("LINE API error: %s, status code: %d", string(body), resp.StatusCode)
	}

	log.Printf("Successfully sent reply to LINE for token: %s", replyToken)
	return nil
}
