package database

import (
	"database/sql"
	"errors"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
)

func (db *appdbimpl) GetMyConversations(userID int) ([]types.Chat, error) {
	var chats []types.Chat

	rows, err := db.c.Query(`
        SELECT id, chat_name, chat_image, is_group, last_msg_id, last_msg_username, last_msg_text, last_msg_time, last_msg_type
        FROM chats
        INNER JOIN user_chats ON id = chat_id
        WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var chat types.Chat
		if err := rows.Scan(&chat.ID, &chat.Name, &chat.Image, &chat.IsGroup, &chat.LastMsgID, &chat.LastMsgUsername, &chat.LastMsgText, &chat.LastMsgTime, &chat.LastMsgType); err != nil {
			return nil, err
		}

		if !chat.IsGroup {
			var user1 int
			var user2 int

			err = db.c.QueryRow(`
			SELECT user1_id, user2_id
			FROM private_chats
			WHERE chat_id = ?`, chat.ID).Scan(&user1, &user2)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}

			if userID == user1 {
				chat.Name, err = db.GetUsernameById(user2)
				if err != nil {
					return nil, err
				}
				chat.Image, err = db.GetImagePath(user2)
				if err != nil {
					return nil, err
				}
			} else {
				chat.Name, err = db.GetUsernameById(user1)
				if err != nil {
					return nil, err
				}
				chat.Image, err = db.GetImagePath(user1)
				if err != nil {
					return nil, err
				}
			}
		}

		var lastRead int
		var unreadCount int

		err = db.c.QueryRow(`
			SELECT message_id
			FROM last_read
			WHERE (user_id, chat_id) = (?,?)`, userID, chat.ID).Scan(&lastRead)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				lastRead = 0
			} else {
				return nil, err
			}
		}

		err = db.c.QueryRow(`
			SELECT COUNT(*)
			FROM messages
			WHERE chat_id = ? AND id > ?`, chat.ID, lastRead).Scan(&unreadCount)
		if err != nil {
			return nil, err
		}

		chat.Unread = unreadCount

		chats = append(chats, chat)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}
