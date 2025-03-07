package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/poprih/ur-monitor/lib/models"
	"github.com/poprih/ur-monitor/lib/repositories"
	"github.com/poprih/ur-monitor/lib/services"
	"github.com/poprih/ur-monitor/pkg/db"
)

// HandleLineWebhook processes incoming LINE webhook events
func HandleLineWebhook(w http.ResponseWriter, r *http.Request, body []byte) error {
	// Parse webhook event
	var webhookEvent models.LineWebhookEvent
	if err := json.Unmarshal(body, &webhookEvent); err != nil {
		return fmt.Errorf("error parsing webhook event: %w", err)
	}

	// Initialize database client
	client, err := db.NewMongoClient()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer client.Disconnect()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(client)
	subRepo := repositories.NewSubscriptionRepository(client)

	// Initialize services
	lineService := services.NewLineService(userRepo, subRepo)

	// Process events
	for _, event := range webhookEvent.Events {
		switch event.Type {
		case "message":
			if event.Message.Type == "text" {
				err := handleTextMessage(event, lineService)
				if err != nil {
					log.Printf("Error processing text message: %v", err)
				}
			}
		case "follow":
			// Handle when a user adds the bot as a friend
			err := handleFollowEvent(event, lineService)
			if err != nil {
				log.Printf("Error handling follow event: %v", err)
			}
		case "unfollow":
			// Handle when a user blocks the bot
			err := handleUnfollowEvent(event, userRepo)
			if err != nil {
				log.Printf("Error handling unfollow event: %v", err)
			}
		default:
			log.Printf("Unhandled event type: %s", event.Type)
		}
	}

	return nil
}

// handleTextMessage processes text messages from users
func handleTextMessage(event models.Event, lineService *services.LineService) error {
	message := event.Message.Text
	userID := event.Source.UserID

	// Check if the user is subscribing to a danchi
	if strings.HasPrefix(strings.ToLower(message), "subscribe") || strings.HasPrefix(message, "登録") {
		// Extract the danchi ID from the message
		parts := strings.Fields(message)
		if len(parts) < 2 {
			return lineService.SendReply(event.ReplyToken, "Please provide a danchi ID to subscribe. Example: 'subscribe 12345'")
		}
		danchiID := parts[1]
		return lineService.SubscribeToDanchi(userID, danchiID, event.ReplyToken)
	} else if strings.HasPrefix(strings.ToLower(message), "unsubscribe") || strings.HasPrefix(message, "解除") {
		// Extract the danchi ID from the message
		parts := strings.Fields(message)
		if len(parts) < 2 {
			return lineService.SendReply(event.ReplyToken, "Please provide a danchi ID to unsubscribe. Example: 'unsubscribe 12345'")
		}
		danchiID := parts[1]
		return lineService.UnsubscribeFromDanchi(userID, danchiID, event.ReplyToken)
	} else if strings.HasPrefix(strings.ToLower(message), "list") || message == "リスト" {
		// List all user subscriptions
		return lineService.ListSubscriptions(userID, event.ReplyToken)
	} else if strings.HasPrefix(strings.ToLower(message), "help") || message == "ヘルプ" {
		// Send help message
		helpMessage := `
Available commands:
- subscribe [danchi_id]: Subscribe to a danchi
- unsubscribe [danchi_id]: Unsubscribe from a danchi  
- list: Show all your subscriptions
- help: Show this help message

Commands also work in Japanese:
- 登録 [danchi_id]: 団地に登録する
- 解除 [danchi_id]: 団地の登録を解除する
- リスト: 登録した団地を表示する
- ヘルプ: このヘルプメッセージを表示する

Example: subscribe 12345
`
		return lineService.SendReply(event.ReplyToken, helpMessage)
	} else {
		// Default response
		return lineService.SendReply(event.ReplyToken, "I don't understand that command. Type 'help' to see available commands.")
	}
}

// handleFollowEvent handles when a user follows (adds) the bot
func handleFollowEvent(event models.Event, lineService *services.LineService) error {
	userID := event.Source.UserID

	// Create or reactivate user
	err := lineService.RegisterUser(userID)
	if err != nil {
		return err
	}

	// Send welcome message
	welcomeMessage := `
Welcome to UR Monitor Bot! 🏠

I can help you monitor UR apartments and notify you when units become available.

How to use:
- subscribe [danchi_id]: Subscribe to a danchi
- unsubscribe [danchi_id]: Unsubscribe from a danchi
- list: Show all your subscriptions
- help: Show all commands

Example: subscribe 12345

日本語も使えます:
- 登録 [danchi_id]: 団地に登録する
- 解除 [danchi_id]: 団地の登録を解除する
- リスト: 登録した団地を表示する
- ヘルプ: コマンド一覧を表示する
`
	return lineService.SendReply(event.ReplyToken, welcomeMessage)
}

// handleUnfollowEvent handles when a user unfollows (blocks) the bot
func handleUnfollowEvent(event models.Event, userRepo *repositories.UserRepository) error {
	userID := event.Source.UserID
	return userRepo.DeactivateUser(userID)
}
