package database

import (
	"math"
)

func (db *appdbimpl) GetLastRead(chatID int) (int, error) {
	var lastRead = math.MaxInt
	var chatMembers []int

	rows, err := db.c.Query(`
		SELECT user_id from user_chats WHERE chat_id = ?`, chatID)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var userId int
		err = rows.Scan(&userId)
		if err != nil {
			return 0, err
		}
		chatMembers = append(chatMembers, userId)
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}

	for _, userId := range chatMembers {
		var userLastRead int

		err = db.c.QueryRow(`
			SELECT message_id FROM last_read WHERE (user_id, chat_id) = (?, ?)`,
			userId, chatID).Scan(&userLastRead)
		if err != nil {
			return 0, err
		}

		if userLastRead < lastRead {
			lastRead = userLastRead
		}

	}

	return lastRead, nil
}
