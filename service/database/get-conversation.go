package database

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
)

func (db *appdbimpl) GetConversation(userId int, chatID int) ([]types.Message, error) {
	var messages []types.Message

	rows, err := db.c.Query(`
        SELECT *
        FROM messages
        WHERE chat_id = ?
        ORDER BY created_at DESC`, chatID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var mostRecentID int

	for rows.Next() {
		var message types.Message
		err = rows.Scan(&message.ID, &message.ChatID, &message.SenderID, &message.Username, &message.Type, &message.Text, &message.MediaURL, &message.IsForward, &message.ReplyTo, &message.CreatedAt)
		if err != nil {
			return nil, err
		}

		if len(messages) == 0 {
			mostRecentID = message.ID
		}

		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	_, err = db.c.Exec(`
		INSERT INTO last_read (user_id, chat_id, message_id)
		VALUES (?, ?, ?)
		ON CONFLICT (user_id, chat_id)
		DO UPDATE SET message_id = excluded.message_id`,
		userId, chatID, mostRecentID)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
