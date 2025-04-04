package models

// LineWebhookEvent represents a LINE webhook event
type LineWebhookEvent struct {
	Destination string  `json:"destination"`
	Events      []Event `json:"events"`
}

// Event represents a LINE event
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
