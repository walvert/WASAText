package types

type User struct {
	ID       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	ImageUrl string `db:"image_url" json:"imageUrl"`
}

type LoginRequest struct {
	Username string `json:"username"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type BearerToken struct {
	Token  string `db:"token"`
	UserID int    `db:"user_id"`
}

type UsernameRequest struct {
	Username string `json:"username"`
}

type SetImageResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	ImageURL string `json:"imageUrl"`
}
