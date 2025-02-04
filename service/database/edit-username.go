package database

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
)

func (db *appdbimpl) EditUsername(id int, username string) (types.User, error) {
	var user types.User

	_, err := db.c.Exec("UPDATE users SET username = ? WHERE id = ?", username, id)
	if err != nil {
		return user, err
	}
	user.Username = username
	user.ID = id
	return user, nil
}
