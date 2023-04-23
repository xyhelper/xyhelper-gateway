package chatresp

import "encoding/json"

type Message struct {
	ID         string   `json:"id"`
	Author     Author   `json:"author"`
	CreateTime float64  `json:"create_time"`
	UpdateTime *string  `json:"update_time"`
	Content    Content  `json:"content"`
	EndTurn    *bool    `json:"end_turn"`
	Weight     float64  `json:"weight"`
	Metadata   Metadata `json:"metadata"`
	Recipient  string   `json:"recipient"`
}

type Author struct {
	Role     string          `json:"role"`
	Name     *string         `json:"name"`
	Metadata map[string]JSON `json:"metadata"`
}

type Content struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type Metadata struct {
	Timestamp   string  `json:"timestamp_"`
	MessageType *string `json:"message_type"`
}

type JSON = json.RawMessage

type ChatCompletion struct {
	Message        Message `json:"message"`
	ConversationId string  `json:"conversation_id"`
	Error          string  `json:"error"`
}
