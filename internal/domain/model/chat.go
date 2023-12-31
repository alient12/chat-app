package model

import "time"

type Chat struct {
	ID        uint64
	People    []uint64
	CreatedAt time.Time
}
