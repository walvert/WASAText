package types

import "time"

type Chat struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	IsGroup     bool      `db:"is_group"`
	LastMsgText string    `db:"last_msg_text"`
	LastMsgTime time.Time `db:"last_msg_time"`
}

type ChatNameRequest struct {
	ChatName string `json:"chat_name"`
}
