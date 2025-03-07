package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/poprih/ur-monitor/lib/models"
	"github.com/poprih/ur-monitor/lib/repositories"
)

const lineReplyAPI = "https://api.line.me/v2/bot/message/reply"

// LineService handles LINE messaging interactions
type LineService struct {
	userRepo *repositories.UserRepository
	subRepo  *repositories.SubscriptionRepository
}

// NewLineService creates a new LINE service instance
func NewLineService(userRepo *repositories.UserRepository, subRepo *repositories.SubscriptionRepository) *LineService {
	return &LineService{
		userRepo: userRepo,
		subRepo:  subRepo,
	}
}

// SendReply sends a reply message to LINE
func (s *LineService) SendReply(replyToken, message string) error {
	// Get LINE channel access token from environment
	accessToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	if accessToken == "" {
		return fmt.Errorf("LINE_CHANNEL_ACCESS_TOKEN is not set")
	}

	// Prepare reply message
	replyMsg := models.LineReplyMessage{
		ReplyToken: replyToken,
		Messages: []models.LineMessage{
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

	return nil
}

// RegisterUser creates or reactivates a user
func (s *LineService) RegisterUser(userID string) error {
	// Check if user exists
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		// Create new user
		newUser := models.User{
			LineUserID: userID,
			Active:     true,
		}
		return s.userRepo.CreateUser(newUser)
	}

	// Reactivate existing user if needed
	if !user.Active {
		return s.userRepo.ActivateUser(userID)
	}

	return nil
}

// SubscribeToDanchi subscribes a user to a danchi
func (s *LineService) SubscribeToDanchi(userID, danchiID, replyToken string) error {
	// Check if the subscription already exists
	exists, err := s.subRepo.SubscriptionExists(userID, danchiID)
	if err != nil {
		return s.SendReply(replyToken, "Error checking subscription status. Please try again later.")
	}

	if exists {
		return s.SendReply(replyToken, fmt.Sprintf("You are already subscribed to danchi %s", danchiID))
	}

	// Create subscription
	subscription := models.Subscription{
		UserID:   userID,
		DanchiID: danchiID,
	}

	err = s.subRepo.AddSubscription(subscription)
	if err != nil {
		return s.SendReply(replyToken, "Failed to add subscription. Please try again later.")
	}

	return s.SendReply(replyToken, fmt.Sprintf("Successfully subscribed to danchi %s! You'll receive notifications when units become available.", danchiID))
}

// UnsubscribeFromDanchi removes a user's subscription to a danchi
func (s *LineService) UnsubscribeFromDanchi(userID, danchiID, replyToken string) error {
	// Check if the subscription exists
	exists, err := s.subRepo.SubscriptionExists(userID, danchiID)
	if err != nil {
		return s.SendReply(replyToken, "Error checking subscription status. Please try again later.")
	}

	if !exists {
		return s.SendReply(replyToken, fmt.Sprintf("You are not subscribed to danchi %s", danchiID))
	}

	// Remove subscription
	err = s.subRepo.RemoveSubscription(userID, danchiID)
	if err != nil {
		return s.SendReply(replyToken, "Failed to remove subscription. Please try again later.")
	}

	return s.SendReply(replyToken, fmt.Sprintf("Successfully unsubscribed from danchi %s.", danchiID))
}

// ListSubscriptions lists all of a user's danchi subscriptions
func (s *LineService) ListSubscriptions(userID, replyToken string) error {
	subscriptions, err := s.subRepo.GetSubscriptionsByUserID(userID)
	if err != nil {
		return s.SendReply(replyToken, "Error retrieving your subscriptions. Please try again later.")
	}

	if len(subscriptions) == 0 {
		return s.SendReply(replyToken, "You don't have any danchi subscriptions yet. Use 'subscribe [danchi_id]' to add one.")
	}

	// Build subscription list message
	message := "Your danchi subscriptions:\n"
	for i, sub := range subscriptions {
		message += fmt.Sprintf("%d. Danchi ID: %s\n", i+1, sub.DanchiID)
	}

	return s.SendReply(replyToken, message)
}
