package model

import "time"

type Message struct {
	ChatID int64
	From   string
	Text   string
	SentAt time.Time
}
