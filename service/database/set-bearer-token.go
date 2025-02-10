package database

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
)

func (db *appdbimpl) UpsertToken(token types.BearerToken) error {
	query := `
        INSERT INTO tokens (user_id, token)
        VALUES (?, ?)
        ON CONFLICT(user_id) DO UPDATE SET token = excluded.token;
    `
	_, err := db.c.Exec(query, token.UserID, token.Token)
	return err
}
