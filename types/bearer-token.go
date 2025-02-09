package types

type BearerToken struct {
	Token  string `db:"token"`
	UserID int    `db:"user_id"`
}
