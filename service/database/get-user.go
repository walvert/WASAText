package database

import (
	"database/sql"
	"errors"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
)

func (db *appdbimpl) GetUser(id int) (types.User, error) {
	var user types.User

	err := db.c.QueryRow("SELECT id, username FROM users WHERE ID = ?", id).
		Scan(&user.ID, &user.Username)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, err // No user found, return nil
		}
		return user, err // Other errors
	}

	return user, nil
}
