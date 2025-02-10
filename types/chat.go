package types

type Chat struct {
	ID      int    `db:"id"`
	Name    string `db:"name"`
	IsGroup bool   `db:"is_group"`
}

type PrivateChat struct {
	User1ID int `db:"user1_id"`
	User2ID int `db:"user2_id"`
	ChatID  int `db:"chat_id"`
}
