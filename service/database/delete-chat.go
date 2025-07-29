package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func (db *appdbimpl) DeleteChat(chatId int) error {
	var isGroup bool
	var imageUrl sql.NullString

	err := db.c.QueryRow(`SELECT is_group, chat_image FROM chats WHERE id = ?`, chatId).Scan(&isGroup, &imageUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("chat not found")
		}
		return err
	}

	if isGroup && imageUrl.Valid && imageUrl.String != "" {
		filePath := filepath.Join("..", "..", imageUrl.String)

		err := os.Remove(filePath)
		if err != nil && !os.IsNotExist(err) {
			fmt.Printf("Warning: failed to delete image file %s: %v\n", filePath, err)
		}
	}

	result, err := db.c.Exec(`DELETE FROM chats WHERE id = ?`, chatId)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("chat not found or unauthorized")
	}

	return nil
}
