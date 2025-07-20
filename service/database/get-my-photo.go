package database

func (db *appdbimpl) GetMyPhoto(token string) (string, error) {
	var imageUrl string
	err := db.c.QueryRow(`
		SELECT image_url
		FROM users u
		INNER JOIN tokens t ON u.id = t.user_id
		WHERE token = ?`, token).Scan(&imageUrl)
	if err != nil {
		return "", err
	}

	if imageUrl == "" {
		return "", ErrEmptyImageURL
	}

	return imageUrl, nil
}
