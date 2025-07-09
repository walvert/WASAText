package database

import (
	"database/sql"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
)

func (db *appdbimpl) GetMyConversations(userID int) ([]types.Chat, error) {
	var chats []types.Chat

	rows, err := db.c.Query(`
        SELECT id, chat_name, is_group, last_msg_text, last_msg_time, last_msg_type
        FROM chats
        INNER JOIN user_chats ON id = chat_id
        WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var chat types.Chat
		if err := rows.Scan(&chat.ID, &chat.Name, &chat.IsGroup, &chat.LastMsgText, &chat.LastMsgTime, &chat.LastMsgType); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, rows.Err()
}
