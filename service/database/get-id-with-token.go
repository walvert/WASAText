package database

func (db *appdbimpl) GetIdWithToken(token string) (int, error) {
	var userId int

	err := db.c.QueryRow("SELECT user_id from tokens WHERE token = ?", token).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
