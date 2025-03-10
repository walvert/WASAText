package types

import "time"

type FirstMessageRequest struct {
	Text      string `json:"text"`
	Type      string `json:"type"`
	MediaURL  string `json:"media_url"`
	ChatName  string `json:"chat_name"`
	Receivers []int  `json:"receivers"`
	IsForward bool   `json:"is_forward"`
}

type MessageRequest struct {
	ChatID    int    `json:"chat_id"`
	Type      string `json:"type"`
	MediaURL  string `json:"media_url"`
	Text      string `json:"text"`
	IsForward bool   `json:"is_forward"`
	ReplyTo   int    `json:"reply_to"`
}

type Message struct {
	ID        int       `db:"id"`
	Type      string    `db:"type"`
	Text      string    `db:"text"`
	ChatID    int       `db:"chat_id"`
	SenderID  int       `db:"sender_id"`
	IsForward bool      `db:"is_forward"`
	ReplyTo   int       `db:"reply_to"`
	CreatedAt time.Time `db:"created_at"`
}

type Comment struct {
	MessageID int `db:"message_id"`
	UserID    int `db:"user_id"`
}

type DeleteRequest struct {
	SenderID int `json:"sender_id"`
}

type ForwardRequest struct {
	Recipients []ForwardRecipient `json:"recipients"`
}

type ForwardRecipient struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

type LastRead struct {
	UserID    int `db:"user_id"`
	ChatID    int `db:"chat_id"`
	MessageID int `db:"message_id"`
}
