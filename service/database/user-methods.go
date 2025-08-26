package database

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
)

func (db *appdbimpl) GetUserById(id int) (types.User, error) {
	var user types.User

	err := db.c.QueryRow("SELECT * FROM users WHERE ID = ?", id).
		Scan(&user)

	return user, err
}

func (db *appdbimpl) GetUsernameById(id int) (string, error) {
	var username string

	err := db.c.QueryRow("SELECT username FROM users WHERE id = ?", id).
		Scan(&username)

	return username, err
}

func (db *appdbimpl) GetUserByUsername(username string) (int, error) {
	var userId int

	err := db.c.QueryRow("SELECT id FROM users WHERE username = ?", username).
		Scan(&userId)

	return userId, err
}

func (db *appdbimpl) CreateUser(username string) (int, error) {
	var userId int

	err := db.c.QueryRow("INSERT INTO users (username) VALUES (?) RETURNING id", username).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (db *appdbimpl) GetUsernameByToken(token string) (string, error) {
	var username string
	err := db.c.QueryRow(`
        SELECT username
        FROM users u
        INNER JOIN tokens t ON u.id = t.user_id
        WHERE token = ?`, token).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

func (db *appdbimpl) GetIdWithToken(token string) (int, error) {
	var userId int

	err := db.c.QueryRow("SELECT user_id from tokens WHERE token = ?", token).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (db *appdbimpl) GetImagePath(userId int) (string, error) {
	var path string
	err := db.c.QueryRow("SELECT image_url FROM users WHERE id = ?", userId).Scan(&path)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (db *appdbimpl) GetUsers(userId int) ([]types.User, error) {
	var users []types.User
	rows, err := db.c.Query("SELECT id, username, image_url FROM users WHERE id != ?", userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user types.User
		err := rows.Scan(&user.ID, &user.Username, &user.ImageUrl)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
