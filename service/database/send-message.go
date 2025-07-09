package database

import "time"

func (db *appdbimpl) SendMessage(chatID int, userID int, text string, msgType string, isForward bool, replyTo int) (int, error) {
	var messageId int
	var timestamp time.Time

	err := db.c.QueryRow(
		"INSERT INTO messages (chat_id, sender_id, text, msg_type, is_forward, reply_to) VALUES (?, ?, ?, ?, ?, ?) RETURNING id, created_at",
		chatID, userID, text, msgType, isForward, replyTo).Scan(&messageId, &timestamp)
	if err != nil {
		return 0, err
	}

	_, err = db.c.Exec(
		`UPDATE chats SET last_msg_text = ?, last_msg_type = ?, last_msg_time = ? WHERE id = ?`,
		text, msgType, timestamp, chatID)

	return messageId, err
}
