package database

import "fmt"

func (db *appdbimpl) DeleteMessage(messageID int) error {
	result, err := db.c.Exec(`DELETE FROM messages
                            WHERE id = ?`,
		messageID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("message not found or unauthorized")
	}

	return nil
}
