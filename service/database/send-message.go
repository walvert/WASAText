package database

func (db *appdbimpl) SendMessage(chatID int, userID int, content string) error {
	_, err := db.c.Exec(
		"INSERT INTO messages (chat_id, sender_id, text) VALUES (?, ?, ?) RETURNING id, chat_id, sender_id, text, created_at",
		chatID, userID, content)

	return err
}
