package model

import "time"

type Message struct {
	ID          uint64             `json:"id,omitempty"`
	ChatID      uint64             `json:"chatid,omitempty"`
	Sender      uint64             `json:"sender,omitempty"`
	Receiver    uint64             `json:"receiver,omitempty"`
	Content     string             `json:"content,omitempty"`
	ContentType MessageContentType `json:"contenttype,omitempty"`
	CreatedAt   time.Time          `json:"createdat,omitempty"`
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
