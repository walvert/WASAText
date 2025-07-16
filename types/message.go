package types

import "time"

type FirstMessageRequest struct {
	Text      string `json:"text"`
	Type      string `json:"type"`
	MediaURL  string `json:"mediaUrl"`
	ChatName  string `json:"chatName"`
	Receivers []int  `json:"receivers"`
	IsForward bool   `json:"isForward"`
}

type MessageRequest struct {
	Type      string `json:"type"`
	MediaURL  string `json:"mediaUrl"`
	Text      string `json:"text"`
	IsForward bool   `json:"isForward"`
	ReplyTo   int    `json:"replyTo"`
}

type Message struct {
	ID        int       `db:"id" json:"id"`
	ChatID    int       `db:"chat_id" json:"chatId"`
	SenderID  int       `db:"sender_id" json:"senderId"`
	Username  string    `db:"username" json:"username"`
	Type      string    `db:"type" json:"type"`
	Text      string    `db:"text" json:"text"`
	MediaURL  string    `db:"media_url" json:"mediaUrl"`
	IsForward bool      `db:"is_forward" json:"isForward"`
	ReplyTo   int       `db:"reply_to" json:"replyTo"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type Comment struct {
	MessageID int `db:"message_id"`
	UserID    int `db:"user_id"`
}

type DeleteRequest struct {
	SenderID int `json:"senderId"`
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
