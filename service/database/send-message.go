package database

func (db *appdbimpl) SendMessage(chatID int, userID int, content string) error {
	_, err := db.c.Exec(
		"INSERT INTO messages (chat_id, user_id, content) VALUES (?, ?, ?) RETURNING id, chat_id, user_id, content, created_at",
		chatID, userID, content)

	return err
}
