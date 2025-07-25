package types

import "time"

type Chat struct {
	ID              int       `db:"id" json:"id"`
	Name            string    `db:"name" json:"name"`
	Image           string    `db:"image" json:"image"`
	IsGroup         bool      `db:"is_group" json:"isGroup"`
	LastMsgID       int       `db:"last_msg_id" json:"lastMsgId"`
	LastMsgUsername string    `db:"last_msg_username" json:"lastMsgUsername"`
	LastMsgText     string    `db:"last_msg_text" json:"lastMsgText"`
	LastMsgTime     time.Time `db:"last_msg_time" json:"lastMsgTime"`
	LastMsgType     string    `db:"last_msg_type" json:"lastMsgType"`
	Unread          int       `json:"unread"`
}

type ChatNameRequest struct {
	ChatName string `json:"chatName"`
}
