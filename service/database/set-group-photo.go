package database

func (db *appdbimpl) SetGroupPhoto(chatId int, imagePath string) error {
	_, err := db.c.Exec("UPDATE chats SET chat_image = ? WHERE id = ?", imagePath, chatId)
	if err != nil {
		return err
	}

	return nil
}
