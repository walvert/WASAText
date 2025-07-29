package database

import "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"

func (db *appdbimpl) GetGroupMembers(chatId int) ([]types.User, error) {
	var members []types.User

	rows, err := db.c.Query(`
		SELECT u.id, u.username, u.image_url
		FROM users u
		INNER JOIN user_chats uc ON u.id = uc.user_id
		WHERE uc.chat_id = ?`, chatId)
	if err != nil {
		return []types.User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var member types.User
		if err := rows.Scan(&member.ID, &member.Username, &member.ImageUrl); err != nil {
			return []types.User{}, err
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return []types.User{}, err
	}

	return members, nil
}
