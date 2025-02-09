package database

import (
	"database/sql"
	"errors"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
)

func (db *appdbimpl) SetToken(token types.BearerToken) error {
	_, err := db.c.Exec("UPDATE tokens SET token = ? WHERE user_id = ?", token.Token, token.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = db.c.Exec("INSERT INTO tokens (token, user_id) VALUES (?, ?)", token.Token, token.UserID)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}
func (db *appdbimpl) UpsertToken(token types.BearerToken) error {
	query := `
        INSERT INTO tokens (user_id, token)
        VALUES (?, ?)
        ON CONFLICT(user_id) DO UPDATE SET token = excluded.token;
    `
	_, err := db.c.Exec(query, token.UserID, token.Token)
	return err
}
