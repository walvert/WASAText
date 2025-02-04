package database

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
)

func (db *appdbimpl) CreateUser(username string) (types.User, error) {
	var user types.User

	err := db.c.QueryRow(
		"INSERT INTO users (username) VALUES ($1) RETURNING id, username",
		username,
	).Scan(&user.ID, &user.Username)

	return user, err
}
