package model

import "time"

type Chat struct {
	ID        uint64    `json:"id,omitempty"`
	People    []uint64  `json:"people,omitempty"`
	CreatedAt time.Time `json:"createdat,omitempty"`
	UpdatedAt time.Time `json:"updatedat,omitempty"`
}

type ChatIDType int

const (
	PrivateChatIDType ChatIDType = iota
	GroupChatIDType
	ChannelChatIDType
	BotChatIDType
	SecretChatIDType
)
