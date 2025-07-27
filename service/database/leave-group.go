package database

import "fmt"

func (db *appdbimpl) LeaveGroup(chatId int, userId int) (chatDeleted bool, err error) {
	res, err := db.c.Exec(`DELETE FROM user_chats WHERE user_id = ? AND chat_id = ?`, userId, chatId)
	if err != nil {
		return false, err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return false, fmt.Errorf("user not found in group")
	}

	_, err = db.c.Exec(`DELETE FROM last_read WHERE user_id = ? AND chat_id = ?`, userId, chatId)
	if err != nil {
	}

	var memberCount int
	err = db.c.QueryRow(`SELECT COUNT(*) FROM user_chats WHERE chat_id = ?`, chatId).Scan(&memberCount)
	if err != nil {
		return false, err
	}

	if memberCount == 0 {
		err = db.DeleteChat(chatId)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}
