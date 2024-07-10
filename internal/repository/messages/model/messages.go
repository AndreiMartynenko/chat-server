package model

import (
	"time"
)

// Message type is the main structure for user message.
type Message struct {
	From      string    `db:"from_user"`
	Text      string    `db:"text"`
	Timestamp time.Time `db:"timestamp"`
}
