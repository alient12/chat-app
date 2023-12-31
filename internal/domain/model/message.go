package model

import "time"

type Message struct {
	ID        uint64
	ChatID    uint64
	Sender    uint64
	Receiver  uint64
	Content   string
	CreatedAt time.Time
}
