package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) EditUserImage(id int, path string) error {
	var existingPath string
	err := db.c.QueryRow("SELECT path FROM image_paths WHERE id = ?", id).Scan(&existingPath)

	if errors.Is(err, sql.ErrNoRows) {
		// No existing entry → Insert new record
		_, err = db.c.Exec("INSERT INTO image_paths (id, path) VALUES (?, ?)", id, path)
		if err != nil {
			return err
		}
	} else if err == nil {
		// Existing entry found → Update the record
		_, err = db.c.Exec("UPDATE image_paths SET path = ? WHERE id = ?", path, id)
		if err != nil {
			return err
		}
	} else {
		return err // Other SQL error
	}

	return nil
}
