package database

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
)

func (db *appdbimpl) CreateChat(chat types.Chat) (types.Chat, error) {
	var myChat types.Chat

	err := db.c.QueryRow(
		"INSERT INTO chats VALUES ($1, $2, $3, $4)",
		chat.ID,
		chat.Name,
		chat.Participants,
		chat.Name).Scan(&myChat)

	return chat, err
}
