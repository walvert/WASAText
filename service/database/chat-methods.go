package database

import (
	"database/sql"
	"errors"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
	"math"
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

func (db *appdbimpl) GetUserChats(userID int) ([]types.Chat, error) {
	var chats []types.Chat

	rows, err := db.c.Query(`
        SELECT id, chat_name, is_group, last_msg_text, last_msg_time, last_msg_type
        FROM chats
        INNER JOIN user_chats ON id = chat_id
        WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var chat types.Chat
		if err := rows.Scan(&chat.ID, &chat.Name, &chat.IsGroup, &chat.LastMsgText, &chat.LastMsgTime, &chat.LastMsgType); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, rows.Err()
}

func (db *appdbimpl) GetConversation(userId int, chatID int) ([]types.Message, error) {
	var messages []types.Message

	rows, err := db.c.Query(`
        SELECT id, chat_id, sender_id, text, created_at, is_forward, reply_to
        FROM messages
        WHERE chat_id = ?
        ORDER BY created_at DESC`, chatID)

	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var mostRecentID int

	for rows.Next() {
		var message types.Message
		err = rows.Scan(&message.ID, &message.ChatID, &message.SenderID, &message.Text, &message.CreatedAt, &message.IsForward, &message.ReplyTo)
		if err != nil {
			return nil, err
		}

		if len(messages) == 0 {
			mostRecentID = message.ID
		}

		messages = append(messages, message)
	}

	_, err = db.c.Exec(`
		INSERT INTO last_read (user_id, chat_id, message_id)
		VALUES (?, ?, ?)
		ON CONFLICT (user_id, chat_id)
		DO UPDATE SET message_id = excluded.message_id`,
		userId, chatID, mostRecentID)

	return messages, rows.Err()
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
