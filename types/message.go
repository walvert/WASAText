package types

import "time"

type FirstMessageRequest struct {
	Text      string `json:"text"`
	ChatName  string `json:"chat_name"`
	Receivers []int  `json:"receivers"`
}

type MessageRequest struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

type Message struct {
	ID        int       `db:"id"`
	Text      string    `db:"text"`
	ChatID    int       `db:"chat_id"`
	SenderID  int       `db:"sender_id"`
	CreatedAt time.Time `db:"created_at"`
}

type Comment struct {
	MessageID int `db:"message_id"`
	UserID    int `db:"user_id"`
}

type DeleteRequest struct {
	SenderID int `json:"sender_id"`
}
