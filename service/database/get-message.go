package database

import "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"

func (db *appdbimpl) GetMessage(messageId int) (types.Message, error) {
	var message types.Message
	err := db.c.QueryRow(`
		SELECT *
		FROM messages
		WHERE id = ?`, messageId).Scan(&message.ID, &message.ChatID, &message.SenderID, &message.Username, &message.Type, &message.Text, &message.MediaURL, &message.IsForward, &message.ReplyTo, &message.CreatedAt)
	if err != nil {
		return types.Message{}, err
	}
	return message, nil
}
