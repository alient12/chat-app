package model

import "time"

type Message struct {
	ID          uint64
	ChatID      uint64
	Sender      uint64
	Receiver    uint64
	Content     string
	ContentType MessageContentType
	CreatedAt   time.Time
}

type MessageContentType int

const (
	TextContentType MessageContentType = iota
	ImageContentType
	VideoContentType
	GifContentType
	FileContentType
	StickerContentType
)
