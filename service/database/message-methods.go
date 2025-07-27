package database

func (db *appdbimpl) GetSenderId(messageId int) (int, error) {
	var senderId int

	err := db.c.QueryRow(`SELECT sender_id FROM messages WHERE id = ?`, messageId).Scan(&senderId)
	if err != nil {
		return 0, err
	}

	return senderId, nil
}

func (db *appdbimpl) GetMessageText(messageID int) (string, error) {
	var text string

	err := db.c.QueryRow(`SELECT text FROM messages WHERE id = ?`,
		messageID).Scan(&text)
	if err != nil {
		return "", err
	}

	return text, nil
}

func (db *appdbimpl) GetMessageType(messageID int) (string, error) {
	var msgType string

	err := db.c.QueryRow(`SELECT type FROM messages WHERE id = ?`,
		messageID).Scan(&msgType)
	if err != nil {
		return "", err
	}

	return msgType, nil
}
