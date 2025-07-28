package database

import "fmt"

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

func (db *appdbimpl) DeleteAllComments(messageID int) error {
	result, err := db.c.Exec(`DELETE FROM comments
                            WHERE message_id = ?`,
		messageID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("comment not found")
	}

	return nil
}
