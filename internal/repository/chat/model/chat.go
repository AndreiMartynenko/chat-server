package model

// Chat type is the main structure for chat.
type Chat struct {
	ID        int64    `db:"id"`
	Usernames []string `db:"usernames"`
}
