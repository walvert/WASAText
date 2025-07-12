package database

func (db *appdbimpl) SetMyPhoto(userId int, path string) error {
	_, err := db.c.Exec("UPDATE users SET image_url = ? WHERE id = ?", path, userId)
	if err != nil {
		return err
	}

	return nil
}
