package database

import (
	"fmt"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
)

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

func (db *appdbimpl) DeletePrivateChat(user1ID int, user2ID int) error {
	if user2ID < user1ID {
		user1ID, user2ID = user2ID, user1ID
	}

	result, err := db.c.Exec(`
        DELETE FROM private_chats
        WHERE user1_id = ? AND user2_id = ?`, user1ID, user2ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("private chat not found")
	}

	return nil
}

func (db *appdbimpl) DeleteUserChat(chatId int) error {
	result, err := db.c.Exec("DELETE FROM user_chats WHERE chat_id = ?", chatId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user chat not found")
	}

	return nil
}

func (db *appdbimpl) UpdateLastMessage(chatId int) error {
	// get last message id
	var highestMessageId int
	err := db.c.QueryRow(`SELECT MAX(id) FROM messages WHERE chat_id = ?`, chatId).Scan(&highestMessageId)
	if err != nil {
		return err
	}

	// get last message info
	var lastMessage types.Message
	lastMessage.ID = highestMessageId
	err = db.c.QueryRow(`
		SELECT sender_id, text, created_at, type
		FROM messages WHERE id = ?`,
		highestMessageId).Scan(&lastMessage.SenderID, &lastMessage.Text, &lastMessage.CreatedAt, &lastMessage.Type)
	if err != nil {
		return err
	}

	// get updated username
	lastMessage.Username, err = db.GetUsernameById(lastMessage.SenderID)
	if err != nil {
		return err
	}

	// update last message info on chat table
	_, err = db.c.Exec(`
		UPDATE chats
		SET last_msg_id = ?, last_msg_username = ?, last_msg_text = ?, last_msg_time = ?, last_msg_type = ?
		WHERE id = ?`,
		lastMessage.ID, lastMessage.Username, lastMessage.Text, lastMessage.CreatedAt, lastMessage.Type, chatId)

	return err
}
