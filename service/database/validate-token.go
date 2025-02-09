package database

import "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"

func (db *appdbimpl) ValidateToken(token types.BearerToken) (bool, error) {
	var dbToken string
	err := db.c.QueryRow("SELECT token from tokens WHERE user_id = ?", token.UserID).Scan(&dbToken)
	if err != nil {
		return false, err
	}
	return dbToken == token.Token, nil
}
