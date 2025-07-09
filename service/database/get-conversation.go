package database

import (
	"database/sql"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
)

func (db *appdbimpl) GetConversation(userId int, chatID int) ([]types.Message, error) {
	var messages []types.Message

	rows, err := db.c.Query(`
        SELECT id, chat_id, sender_id, text, created_at, is_forward, reply_to
        FROM messages
        WHERE chat_id = ?
        ORDER BY created_at DESC`, chatID)

	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var mostRecentID int

	for rows.Next() {
		var message types.Message
		err = rows.Scan(&message.ID, &message.ChatID, &message.SenderID, &message.Text, &message.CreatedAt, &message.IsForward, &message.ReplyTo)
		if err != nil {
			return nil, err
		}

		if len(messages) == 0 {
			mostRecentID = message.ID
		}

		messages = append(messages, message)
	}

	_, err = db.c.Exec(`
		INSERT INTO last_read (user_id, chat_id, message_id)
		VALUES (?, ?, ?)
		ON CONFLICT (user_id, chat_id)
		DO UPDATE SET message_id = excluded.message_id`,
		userId, chatID, mostRecentID)

	return messages, rows.Err()
}
