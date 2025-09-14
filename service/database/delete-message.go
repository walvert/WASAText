package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) DeleteMessage(messageId int) (chatDeleted bool, err error) {
	var chatId int

	// get chat id
	err = db.c.QueryRow(`SELECT chat_id FROM messages WHERE id = ?`, messageId).Scan(&chatId)
	if err != nil {
		return false, fmt.Errorf("message not found: %w", err)
	}

	// delete message
	result, err := db.c.Exec(`DELETE FROM messages WHERE id = ?`, messageId)
	if err != nil {
		return false, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return false, fmt.Errorf("message not found or unauthorized")
	}

	// delete comments if any
	err = db.DeleteAllComments(messageId)
	if err != nil {
	}

	// check if the message deleted was the only one in the chat
	var count int
	err = db.c.QueryRow(`SELECT COUNT(*) FROM messages WHERE chat_id = ?`, chatId).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		var user1ID, user2ID int
		err = db.c.QueryRow(`SELECT user1_id, user2_id FROM private_chats WHERE chat_id = ?`, chatId).Scan(&user1ID, &user2ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return false, err
		}

		if err == nil {
			err = db.DeletePrivateChat(user1ID, user2ID)
			if err != nil {
				return false, err
			}
		}

		err = db.DeleteChat(chatId)
		if err != nil {
			return false, err
		}

		err = db.DeleteUserChat(chatId)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	// check if the message deleted was the most recent in the chat
	var chatLastMessageId int
	err = db.c.QueryRow(`SELECT last_msg_id FROM chats WHERE id = ?`, chatId).Scan(&chatLastMessageId)
	if err != nil {
		return false, err
	}
	if chatLastMessageId == messageId {
		err = db.UpdateLastMessage(chatId)
	}

	return false, nil
}
