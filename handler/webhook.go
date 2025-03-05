package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

const lineAPI = "https://api.line.me/v2/bot/message/reply"

// LINE Webhook Request Structs
type LineWebhookRequest struct {
	Events []struct {
		ReplyToken string `json:"replyToken"`
		Type       string `json:"type"`
		Message    struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"events"`
}

type ReplyMessage struct {
	ReplyToken string `json:"replyToken"`
	Messages   []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"messages"`
}

// Handle incoming LINE messages
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	var req LineWebhookRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	for _, event := range req.Events {
		if event.Type == "message" && event.Message.Type == "text" {
			replyText(event.ReplyToken, "你发送了: "+event.Message.Text)
		}
	}

	w.WriteHeader(http.StatusOK)
}

// Reply to the user
func replyText(replyToken, message string) {
	accessToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	if accessToken == "" {
		log.Println("Missing LINE_CHANNEL_ACCESS_TOKEN")
		return
	}

	reply := ReplyMessage{
		ReplyToken: replyToken,
		Messages: []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		}{{Type: "text", Text: message}},
	}

	data, _ := json.Marshal(reply)
	req, _ := http.NewRequest("POST", lineAPI, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending reply:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("Reply sent successfully")
}

func main() {
	http.HandleFunc("/webhook", webhookHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server started on port:", port)
	http.ListenAndServe(":"+port, nil)
}
