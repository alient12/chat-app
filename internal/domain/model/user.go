package model

type User struct {
	ID        uint64 `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	Image     string `json:"image,omitempty"`
	Bio       string `json:"bio,omitempty"`
}

type IDType int

const (
	UserIDType IDType = iota
	GroupIDType
	ChannelIDType
	BotIDType
)
