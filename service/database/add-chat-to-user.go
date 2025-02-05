package database

func (db *appdbimpl) AddChatToUser(chatID int, userID int) error {
	_, err := db.c.Exec("INSERT INTO user_chats (user_id, chat_id) VALUES (?, ?)", userID, chatID)
	return err
}
