package types

type Chat struct {
	ID           int       `json:"id"`
	Participants []User    `json:"participants"`
	Messages     []Message `json:"messages"`
}
