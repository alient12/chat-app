package model

type Profile struct {
	UserID    uint64
	Bio       string
	Avatar    []File
	AvatarLow []byte
	Settings  map[string]string
}
