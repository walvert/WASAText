package database

import (
	"fmt"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
)

func (db *appdbimpl) GetUserById(id int) (types.User, error) {
	var user types.User

	err := db.c.QueryRow("SELECT * FROM users WHERE ID = ?", id).
		Scan(&user)

	return user, err
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

func (db *appdbimpl) SetGroupName(chatId int, chatName string) error {
	_, err := db.c.Exec("UPDATE chats SET chat_name = ? WHERE id = ?", chatName, chatId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) LeaveGroup(chatId int, userId int) error {
	result, err := db.c.Exec(`DELETE FROM user_chats
                            WHERE user_id = ? AND chat_id = ?`,
		userId, chatId)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
