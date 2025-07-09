package database

func (db *appdbimpl) CreateChat(chatName string, isGroup bool) (int, error) {
	var chatID int
	err := db.c.QueryRow(
		"INSERT INTO chats (chat_name, is_group) VALUES (?,?) RETURNING id",
		chatName,
		isGroup).Scan(&chatID)

	return chatID, err
}
