package database

import "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"

func (db *appdbimpl) SetMyUsername(user types.User) error {
	_, err := db.c.Exec("UPDATE users SET username = ? WHERE id = ?", user.Username, user.ID)
	if err != nil {
		return err
	}
	return nil
}
