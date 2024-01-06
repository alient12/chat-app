package handler

import (
	"chatapp/internal/domain/repository/chatrepo"
	"chatapp/internal/domain/repository/messagerepo"
	"chatapp/internal/domain/repository/userrepo"
	"math/rand"
)

type Message struct {
	urepo  userrepo.Repository
	chrepo chatrepo.Repository
	repo   messagerepo.Repository
}

func NewMessage(repo messagerepo.Repository, urepo userrepo.Repository, chrepo chatrepo.Repository) *Message {
	return &Message{
		urepo:  urepo,
		chrepo: chrepo,
		repo:   repo,
	}
}

var (
	message_count uint64
)

func GenerateMessageID() uint64 {
	id := (message_count << 32) | (uint64(rand.Uint32()))
	return id
}
