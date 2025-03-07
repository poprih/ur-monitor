package models

// User represents a LINE bot user
type User struct {
	LineUserID string `bson:"line_user_id"`
	Active     bool   `bson:"active"`
}

// Subscription represents a user's subscription to a danchi
type Subscription struct {
	UserID   string `bson:"user_id"`
	DanchiID string `bson:"danchi_id"`
}

// Property represents a UR property unit
type Property struct {
	Name       string `json:"name"`
	FloorType  string `json:"floor"`
	RentNormal string `json:"rent_normal"`
	DetailLink string `json:"roomDetailLink"`
	DanchiID   string `json:"danchi_id"`
}

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

// LineReplyMessage represents a LINE reply message
type LineReplyMessage struct {
	ReplyToken string        `json:"replyToken"`
	Messages   []LineMessage `json:"messages"`
}

// LineMessage represents a LINE message
type LineMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// LineNotifyMessage represents a LINE Notify message
type LineNotifyMessage struct {
	Message string `json:"message"`
}
