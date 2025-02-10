package database

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
)

func (db *appdbimpl) CreateChat(chatName string, isGroup bool) (int, error) {
	var chatID int
	err := db.c.QueryRow(
		"INSERT INTO chats (chat_name, is_group) VALUES (?,?) RETURNING id",
		chatName,
		isGroup).Scan(&chatID)

	return chatID, err
}

func (db *appdbimpl) AddChatToUser(userID int, chatID int) error {
	_, err := db.c.Exec("INSERT INTO user_chats (user_id, chat_id) VALUES (?, ?)", userID, chatID)
	return err
}
func (db *appdbimpl) AddPrivateChat(chat types.PrivateChat) error {
	_, err := db.c.Exec(
		"INSERT INTO private_chats (user1_id, user2_id, chat_id) VALUES (?, ?, ?)",
		chat.User1ID,
		chat.User2ID,
		chat.ChatID)
	return err
}
func (db *appdbimpl) GetPrivateChatID(user1ID int, user2ID int) (int, error) {
	var chatID int
	err := db.c.QueryRow("SELECT chat_id FROM private_chats WHERE (user1_id = ? AND user2_id = ?)",
		user1ID,
		user2ID).Scan(&chatID)
	return chatID, err
}
