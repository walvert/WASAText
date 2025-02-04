package types

import "time"

type Message struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	ChatID    int       `json:"chat_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
