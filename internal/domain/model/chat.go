package model

import "time"

type Chat struct {
	ID        uint64
	People    []uint64
	CreatedAt time.Time
}

type ChatIDType int

const (
	PrivateChatIDType ChatIDType = iota
	GroupChatIDType
	ChannelChatIDType
	BotChatIDType
	SecretChatIDType
)
