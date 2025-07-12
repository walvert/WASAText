package database

func (db *appdbimpl) SetGroupName(chatId int, chatName string) error {
	_, err := db.c.Exec("UPDATE chats SET chat_name = ? WHERE id = ?", chatName, chatId)
	if err != nil {
		return err
	}
	return nil
}
