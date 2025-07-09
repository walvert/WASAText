package database

import "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"

func (db *appdbimpl) GetChatInfo(chatId int) (types.Chat, error) {
	var chat types.Chat

	err := db.c.QueryRow(`
        SELECT *
        FROM chats
        WHERE id = ?`, chatId).Scan(&chat.ID, &chat.Name, &chat.IsGroup, &chat.LastMsgText, &chat.LastMsgTime, &chat.LastMsgType)

	return chat, err
}
