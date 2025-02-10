package database

import (
	"database/sql"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
)

func (db *appdbimpl) GetConversation(chatID int) ([]types.Message, error) {
	var messages []types.Message

	rows, err := db.c.Query(`
        SELECT id, chat_id, user_id, content, created_at
        FROM messages
        WHERE chat_id = ?`, chatID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)

	for rows.Next() {
		var message types.Message
		if err := rows.Scan(&message.ID, &message.ChatID, &message.SenderID, &message.Text, &message.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, rows.Err()
}
