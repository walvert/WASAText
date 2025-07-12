package database

func (db *appdbimpl) SetLastRead(userId int, chatId int) error {
	_, err := db.c.Exec(`
		INSERT INTO last_read (user_id, chat_id, message_id)
		VALUES (?, ?, ?)
		ON CONFLICT (user_id, chat_id)
		DO UPDATE SET message_id = excluded.message_id`,
		userId, chatId, 0)

	return err
}
