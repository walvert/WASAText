package database

import (
	"log"
)

func (db *appdbimpl) UsernameExists(username string) bool {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
	if err != nil {
		log.Println("Error checking username:", err)
		return false
	}
	return exists
}
