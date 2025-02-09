package database

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
)

func (db *appdbimpl) CreateChat(chat types.Chat) error {
	_, err := db.c.Exec(
		"INSERT INTO chats (id, chat_name, is_group) VALUES (?,?,?)",
		chat.ID,
		chat.Name,
		chat.IsGroup)

	return err
}
