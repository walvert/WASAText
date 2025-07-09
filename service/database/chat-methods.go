package database

import (
	"errors"
	"math"
)

func (db *appdbimpl) AddChatToUser(userID int, chatID int) error {
	_, err := db.c.Exec("INSERT INTO user_chats (user_id, chat_id) VALUES (?, ?)", userID, chatID)
	return err
}

func (db *appdbimpl) AddToGroup(chatID int, userID int) error {
	var isGroup bool

	err := db.c.QueryRow("SELECT is_group FROM chats WHERE id = ?", chatID).Scan(&isGroup)
	if err != nil {
		return err
	}

	if !isGroup {
		return errors.New("unauthorized")
	}

	_, err = db.c.Exec("INSERT INTO user_chats (user_id, chat_id) VALUES (?, ?)", userID, chatID)
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

func (db *appdbimpl) GetLastRead(chatID int) (int, error) {
	var lastRead = math.MaxInt
	var chatMembers []int

	rows, err := db.c.Query(`
		SELECT user_id from user_chats WHERE chat_id = ?`, chatID)
	if err != nil {
		return 0, err
	}

	for rows.Next() {
		var userId int
		err = rows.Scan(&userId)
		if err != nil {
			return 0, err
		}
		chatMembers = append(chatMembers, userId)
	}
	err = rows.Close()
	if err != nil {
		return 0, err
	}

	for _, userId := range chatMembers {
		var userLastRead int

		err = db.c.QueryRow(`
			SELECT message_id from last_read WHERE chat_id = ? AND user_id = ?`,
			chatID, userId).Scan(&userLastRead)
		if err != nil {
			return 0, err
		}

		if userLastRead < lastRead {
			lastRead = userLastRead
		}

	}

	return lastRead, nil
}
