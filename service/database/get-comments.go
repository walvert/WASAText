package database

func (db *appdbimpl) GetComments(messageId int) ([]string, error) {
	var comments []string

	rows, err := db.c.Query(`
		SELECT user_id
		FROM comments
		WHERE message_id = ?`, messageId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		username, err := db.GetUsernameById(id)
		if err != nil {
			return nil, err
		}
		comments = append(comments, username)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
