package database

import "fmt"

func (db *appdbimpl) LeaveGroup(chatId int, userId int) error {
	result, err := db.c.Exec(`DELETE FROM user_chats
                            WHERE (user_id, chat_id) = (?,?)`,
		userId, chatId)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
