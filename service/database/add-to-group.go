package database

import (
	"errors"
)

func (db *appdbimpl) AddToGroup(chatID int, userID int) error {
	var isGroup bool
	var exists bool

	err := db.c.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM user_chats WHERE user_id = ? AND chat_id = ?)",
		userID, chatID,
	).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return ErrAlreadyExists
	}

	err = db.c.QueryRow("SELECT is_group FROM chats WHERE id = ?", chatID).Scan(&isGroup)
	if err != nil {
		return err
	}

	if !isGroup {
		return errors.New("bad request")
	}

	_, err = db.c.Exec("INSERT INTO user_chats (user_id, chat_id) VALUES (?, ?)", userID, chatID)
	return err
}
