package database

import "fmt"

func (db *appdbimpl) LeaveGroup(chatId int, userId int) error {
	res, err := db.c.Exec(`DELETE FROM user_chats WHERE user_id = ? AND chat_id = ?`, userId, chatId)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found in group")
	}

	_, err = db.c.Exec(`DELETE FROM last_read WHERE user_id = ? AND chat_id = ?`, userId, chatId)
	if err != nil {
	}

	return nil
}
