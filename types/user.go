package types

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type LoginRequest struct {
	Username string `json:"username"`
}

type LoginResponse struct {
	Token string `json:"token"`
	ID    int    `json:"id"`
}

type BearerToken struct {
	Token  string `db:"token"`
	UserID int    `db:"user_id"`
}
