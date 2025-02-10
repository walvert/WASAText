package database

import (
	"database/sql"
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

func (db *appdbimpl) GetUserChats(userID int) ([]types.Chat, error) {
	var chats []types.Chat

	rows, err := db.c.Query(`
        SELECT id, chat_name, is_group
        FROM chats
        INNER JOIN user_chats ON id = chat_id
        WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)

	for rows.Next() {
		var chat types.Chat
		if err := rows.Scan(&chat.ID, &chat.Name, &chat.IsGroup); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, rows.Err()
}
