package database

import "fmt"

func (db *appdbimpl) DeleteChat(chatId int) error {
	result, err := db.c.Exec(`DELETE FROM chats
                            WHERE id = ?`,
		chatId)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("message not found or unauthorized")
	}

	return nil
}
