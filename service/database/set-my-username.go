package database

import (
	"database/sql"
	"errors"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
)

func (db *appdbimpl) SetMyUserName(user types.User) error {
	var existingUserId int
	err := db.c.QueryRow("SELECT id FROM users WHERE username = ? AND id != ?",
		user.Username, user.ID).Scan(&existingUserId)

	if err == nil {
		return ErrUsernameAlreadyExists
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	_, err = db.c.Exec("UPDATE users SET username = ? WHERE id = ?", user.Username, user.ID)
	if err != nil {
		return err
	}
	return nil
}
