package types

type Chat struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Participants []User    `json:"participants"`
	Messages     []Message `json:"messages"`
}
