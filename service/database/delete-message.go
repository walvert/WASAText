package database

import "fmt"

func (db *appdbimpl) DeleteMessage(messageId int) (chatDeleted bool, err error) {
	var chatId int
	err = db.c.QueryRow(`SELECT chat_id FROM messages WHERE id = ?`, messageId).Scan(&chatId)
	if err != nil {
		return false, fmt.Errorf("message not found")
	}

	result, err := db.c.Exec(`DELETE FROM messages WHERE id = ?`, messageId)
	if err != nil {
		return false, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return false, fmt.Errorf("message not found or unauthorized")
	}

	var count int
	err = db.c.QueryRow(`SELECT COUNT(*) FROM messages WHERE chat_id = ?`, chatId).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		err = db.DeleteChat(chatId)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}
