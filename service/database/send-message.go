package database

import "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"

func (db *appdbimpl) SendMessage(message types.Message) (types.Message, error) {
	var myMessage types.Message

	err := db.c.QueryRow(
		"INSERT INTO messages (id, chat_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?)",
		message.ID,
		message.ChatID,
		message.UserID,
		message.Text,
		message.CreatedAt).Scan(&myMessage)

	return myMessage, err
}
