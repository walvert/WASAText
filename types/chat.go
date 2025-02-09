package types

type Chat struct {
	ID      int    `db:"id"`
	Name    string `db:"name"`
	IsGroup bool   `db:"is_group"`
}
