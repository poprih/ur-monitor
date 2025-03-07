package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/poprih/ur-monitor/lib/config"
	"github.com/poprih/ur-monitor/lib/models"
	"github.com/poprih/ur-monitor/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func LineWebhook(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	cfg, err := config.GetConfig()
	if err != nil {
		log.Printf("Error loading config: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Initialize LINE bot client
	bot, err := linebot.New(cfg.LineChannelSecret, cfg.LineChannelToken)
	if err != nil {
		log.Printf("Error initializing LINE bot: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Parse incoming webhook request
	events, err := bot.ParseRequest(r)
	if err != nil {
		log.Printf("Error parsing webhook request: %v", err)
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Process each event
	for _, event := range events {
		if err := handleLineEvent(bot, event); err != nil {
			log.Printf("Error handling event: %v", err)
		}
	}

	w.WriteHeader(http.StatusOK)
}

// handleLineEvent processes a single LINE event
func handleLineEvent(bot *linebot.Client, event *linebot.Event) error {
	userID := event.Source.UserID

	// Handle different event types
	switch event.Type {
	case linebot.EventTypeFollow:
		// User followed the bot - register them
		if err := registerUser(userID); err != nil {
			return err
		}

		// Send welcome message
		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(
			"UR物件モニターへようこそ！\n"+
				"「登録」と入力して、希望の団地を登録しましょう。"),
		).Do(); err != nil {
			return err
		}

	case linebot.EventTypeUnfollow:
		// User unfollowed the bot - unregister them
		if err := unregisterUser(userID); err != nil {
			return err
		}

	case linebot.EventTypeMessage:
		// Handle message events
		switch message := event.Message.(type) {
		case *linebot.TextMessage:
			return handleTextMessage(bot, event, message)
		}
	}

	return nil
}

// handleTextMessage processes a text message from a user
func handleTextMessage(bot *linebot.Client, event *linebot.Event, message *linebot.TextMessage) error {
	text := message.Text
	userID := event.Source.UserID

	// Handle different command types
	switch text {
	case "登録":
		// List available danchi for registration
		cfg, err := config.GetConfig()
		if err != nil {
			return err
		}

		danchiButtons := []linebot.FlexComponent{}
		for _, danchi := range cfg.DanchiList {
			danchiButtons = append(danchiButtons, linebot.NewButton(
				"登録: "+danchi,
				linebot.NewPostbackAction("選択", "register:"+danchi, "", ""),
			))
		}

		container := linebot.NewBubbleContainer(
			linebot.NewBoxComponent(
				linebot.FlexComponentTypeText,
				linebot.FlexComponentTypeBox,
				&linebot.BoxComponentBubble{
					Body: &linebot.BoxComponent{
						Type:     linebot.FlexComponentTypeBox,
						Layout:   linebot.FlexBoxLayoutTypeVertical,
						Contents: danchiButtons,
					},
				},
			),
		)

		if _, err := bot.ReplyMessage(
			event.ReplyToken,
			linebot.NewFlexMessage("団地を選択してください", container),
		).Do(); err != nil {
			return err
		}

	case "確認":
		// Show current subscriptions
		subs, err := getUserSubscriptions(userID)
		if err != nil {
			return err
		}

		if len(subs) == 0 {
			if _, err := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage("登録されている団地はありません。「登録」と入力して団地を登録しましょう。"),
			).Do(); err != nil {
				return err
			}
		} else {
			subscriptionText := "登録中の団地:\n"
			for _, sub := range subs {
				subscriptionText += "- " + sub.DanchiName + "\n"
			}

			if _, err := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage(subscriptionText),
			).Do(); err != nil {
				return err
			}
		}

	case "解除":
		// List subscriptions for unregistration
		subs, err := getUserSubscriptions(userID)
		if err != nil {
			return err
		}

		if len(subs) == 0 {
			if _, err := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage("登録されている団地はありません。"),
			).Do(); err != nil {
				return err
			}
		} else {
			danchiButtons := []linebot.FlexComponent{}
			for _, sub := range subs {
				danchiButtons = append(danchiButtons, linebot.NewButton(
					"解除: "+sub.DanchiName,
					linebot.NewPostbackAction("選択", "unregister:"+sub.DanchiName, "", ""),
				))
			}

			container := linebot.NewBubbleContainer(
				linebot.NewBoxComponent(
					linebot.FlexComponentTypeText,
					linebot.FlexComponentTypeBox,
					&linebot.BoxComponentBubble{
						Body: &linebot.BoxComponent{
							Type:     linebot.FlexComponentTypeBox,
							Layout:   linebot.FlexBoxLayoutTypeVertical,
							Contents: danchiButtons,
						},
					},
				),
			)

			if _, err := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewFlexMessage("解除する団地を選択してください", container),
			).Do(); err != nil {
				return err
			}
		}

	case "ヘルプ":
		// Show help message
		helpMessage := "使い方:\n" +
			"「登録」 - 希望の団地を登録\n" +
			"「確認」 - 登録中の団地を確認\n" +
			"「解除」 - 登録を解除\n" +
			"「ヘルプ」 - このヘルプを表示"

		if _, err := bot.ReplyMessage(
			event.ReplyToken,
			linebot.NewTextMessage(helpMessage),
		).Do(); err != nil {
			return err
		}

	default:
		// Unknown command
		if _, err := bot.ReplyMessage(
			event.ReplyToken,
			linebot.NewTextMessage("コマンドが認識できません。「ヘルプ」と入力してコマンド一覧を確認してください。"),
		).Do(); err != nil {
			return err
		}
	}

	return nil
}

// registerUser adds a new user to the database
func registerUser(userID string) error {
	database, err := db.GetDatabase()
	if err != nil {
		return err
	}

	collection := database.Collection("users")
	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"userId": userID},
		bson.M{"$setOnInsert": bson.M{
			"userId":    userID,
			"createdAt": time.Now(),
		}},
		options.Update().SetUpsert(true),
	)

	return err
}

// unregisterUser removes a user and their subscriptions
func unregisterUser(userID string) error {
	database, err := db.GetDatabase()
	if err != nil {
		return err
	}

	// Remove user
	userCollection := database.Collection("users")
	_, err = userCollection.DeleteOne(context.Background(), bson.M{"userId": userID})
	if err != nil {
		return err
	}

	// Remove all subscriptions
	subCollection := database.Collection("subscriptions")
	_, err = subCollection.DeleteMany(context.Background(), bson.M{"userId": userID})

	return err
}

// getUserSubscriptions retrieves all danchi subscriptions for a user
func getUserSubscriptions(userID string) ([]models.Subscription, error) {
	database, err := db.GetDatabase()
	if err != nil {
		return nil, err
	}

	collection := database.Collection("subscriptions")
	cursor, err := collection.Find(
		context.Background(),
		bson.M{"userId": userID},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var subscriptions []models.Subscription
	if err = cursor.All(context.Background(), &subscriptions); err != nil {
		return nil, err
	}

	return subscriptions, nil
}
