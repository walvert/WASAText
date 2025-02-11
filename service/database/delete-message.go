package database

import "fmt"

func (db *appdbimpl) DeleteMessage(messageId int) error {
	result, err := db.c.Exec(`DELETE FROM messages
                            WHERE id = ?`,
		messageId)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("message not found or unauthorized")
	}

	return nil
}

func (db *appdbimpl) LeaveGroup(chatId int, userId int) error {
	result, err := db.c.Exec(`DELETE FROM user_chats
                            WHERE user_id = ? AND chat_id = ?`,
		userId, chatId)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("message not found or unauthorized")
	}

	return nil
}

func (db *appdbimpl) GetSenderId(messageId int) (int, error) {
	var senderId int

	err := db.c.QueryRow(`SELECT sender_id FROM messages WHERE id = ?`, messageId).Scan(&senderId)
	if err != nil {
		return 0, err
	}

	return senderId, nil
}
