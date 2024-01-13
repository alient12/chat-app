package model

import "time"

type File struct {
	ID               uint64
	UserID           uint64
	FileName         string
	Size             int64
	ContentType      string
	FilePath         string
	ChatIDs          []uint64
	Metadata         map[string]string
	IsProfileContent bool
	CreatedAt        time.Time
}
