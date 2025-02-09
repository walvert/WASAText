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

func (db *appdbimpl) GetUserByUsername(username string) (types.User, error) {
	var user types.User

	err := db.c.QueryRow("SELECT * FROM users WHERE username = ?", username).
		Scan(&user.ID, &user.Username)

	return user, err
}

func (db *appdbimpl) CreateUser(username string) (types.User, error) {
	var user types.User

	err := db.c.QueryRow("INSERT INTO users (username) VALUES ($1) RETURNING id, username", username).Scan(&user.ID, &user.Username)
	if err != nil {
		return user, err
	}
	return user, nil
}

/*func (db *appdbimpl) SetMyPhoto(id int, path string) error {
	var existingPath string
	err := db.c.QueryRow("SELECT path FROM image_paths WHERE id = ?", id).Scan(&existingPath)

	if errors.Is(err, sql.ErrNoRows) {
		// No existing entry → Insert new record
		_, err = db.c.Exec("INSERT INTO image_paths (id, path) VALUES (?, ?)", id, path)
		if err != nil {
			return err
		}
	} else if err == nil {
		// Existing entry found → Update the record
		_, err = db.c.Exec("UPDATE image_paths SET path = ? WHERE id = ?", path, id)
		if err != nil {
			return err
		}
	} else {
		return err // Other SQL error
	}

	return nil
}*/

func (db *appdbimpl) SetMyUsername(user types.User) error {
	_, err := db.c.Exec("UPDATE users SET username = ? WHERE id = ?", user.Username, user.ID)
	if err != nil {
		return err
	}
	return nil
}
