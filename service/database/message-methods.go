package database

import (
	"fmt"
)

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

func (db *appdbimpl) GetSenderId(messageId int) (int, error) {
	var senderId int

	err := db.c.QueryRow(`SELECT sender_id FROM messages WHERE id = ?`, messageId).Scan(&senderId)
	if err != nil {
		return 0, err
	}

	return senderId, nil
}

func (db *appdbimpl) DeleteComment(messageID int, userID int) error {
	result, err := db.c.Exec(`DELETE FROM comments
                            WHERE message_id = ? AND user_id = ?`,
		messageID, userID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("comment not found")
	}

	return nil
}

func (db *appdbimpl) CommentMessage(messageID int, userID int) error {
	_, err := db.c.Exec(
		"INSERT INTO comments (message_id, user_id) VALUES (?, ?)",
		messageID, userID)

	return err
}

func (db *appdbimpl) GetMessageText(messageID int) (string, error) {
	var text string

	err := db.c.QueryRow(`SELECT text FROM messages WHERE id = ?`,
		messageID).Scan(&text)
	if err != nil {
		return "", err
	}

	return text, nil
}

func (db *appdbimpl) GetMessageType(messageID int) (string, error) {
	var msgType string

	err := db.c.QueryRow(`SELECT msg_type FROM messages WHERE id = ?`,
		messageID).Scan(&msgType)
	if err != nil {
		return "", err
	}

	return msgType, nil
}
