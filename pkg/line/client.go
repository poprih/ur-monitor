package line

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// LineClient represents a LINE messaging API client
type LineClient struct {
	channelToken string
	httpClient   *http.Client
}

// NewClient creates a new LINE client
func NewLineClient(channelToken string) *LineClient {
	return &LineClient{
		channelToken: channelToken,
		httpClient:   &http.Client{},
	}
}

// SendPushMessage sends a push message to a LINE user
func (c *LineClient) SendPushMessage(userID, message string) error {
	payload := map[string]interface{}{
		"to": userID,
		"messages": []map[string]string{
			{
				"type": "text",
				"text": message,
			},
		},
	}

	return c.sendRequest("https://api.line.me/v2/bot/message/push", payload)
}

// SendReplyMessage sends a reply message to a LINE user
func (c *LineClient) SendReplyMessage(replyToken, message string) error {
	payload := map[string]interface{}{
		"replyToken": replyToken,
		"messages": []map[string]string{
			{
				"type": "text",
				"text": message,
			},
		},
	}

	return c.sendRequest("https://api.line.me/v2/bot/message/reply", payload)
}

// sendRequest sends a request to the LINE API
func (c *LineClient) sendRequest(url string, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.channelToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("LINE API error: %s (status code: %d)", string(body), resp.StatusCode)
	}

	return nil
} 
