package database

import "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"

func (db *appdbimpl) GetForwardInfo(messageId int) (types.ForwardInfo, error) {
	var info types.ForwardInfo
	err := db.c.QueryRow(`SELECT type, text, media_url FROM messages WHERE id = ?`,
		messageId).Scan(&info.Type, &info.Text, &info.MediaUrl)
	if err != nil {
		return info, err
	}
	return info, nil
}
