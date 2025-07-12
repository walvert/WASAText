package database

func (db *appdbimpl) AddChatToUser(userID int, chatID int) error {
	_, err := db.c.Exec("INSERT INTO user_chats (user_id, chat_id) VALUES (?, ?)", userID, chatID)
	return err
}

func (db *appdbimpl) AddPrivateChat(user1Id int, user2Id int, chatId int) error {
	_, err := db.c.Exec(
		"INSERT INTO private_chats (user1_id, user2_id, chat_id) VALUES (?, ?, ?)",
		user1Id,
		user2Id,
		chatId)
	return err
}

func (db *appdbimpl) GetPrivateChatID(user1ID int, user2ID int) (int, error) {
	var chatID int
	err := db.c.QueryRow("SELECT chat_id FROM private_chats WHERE (user1_id = ? AND user2_id = ?)",
		user1ID,
		user2ID).Scan(&chatID)
	return chatID, err
}
