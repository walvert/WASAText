package database

func (db *appdbimpl) CommentMessage(messageID int, userID int) error {
	_, err := db.c.Exec(
		"INSERT INTO comments (message_id, user_id) VALUES (?, ?)",
		messageID, userID)

	return err
}
