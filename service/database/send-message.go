package database

import "time"

func (db *appdbimpl) SendMessage(chatID int, userID int, username string, msgType string, text string, mediaUrl string, isForward bool, replyTo int) (int, error) {
	var messageId int
	var timestamp time.Time

	err := db.c.QueryRow(
		"INSERT INTO messages (chat_id, sender_id, username, type, text, media_url, is_forward, reply_to) VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING id, created_at",
		chatID, userID, username, msgType, text, mediaUrl, isForward, replyTo).Scan(&messageId, &timestamp)
	if err != nil {
		return 0, err
	}

	_, err = db.c.Exec(
		`UPDATE chats SET last_msg_id = ?, last_msg_username = ?, last_msg_text = ?, last_msg_type = ?, last_msg_time = ? WHERE id = ?`,
		messageId, username, text, msgType, timestamp, chatID)
	if err != nil {
		return 0, err
	}

	_, err = db.c.Exec(`
		INSERT INTO last_read (user_id, chat_id, message_id)
		VALUES (?, ?, ?)
		ON CONFLICT (user_id, chat_id)
		DO UPDATE SET message_id = excluded.message_id`,
		userID, chatID, messageId)
	if err != nil {
		return 0, err
	}

	return messageId, nil
}
