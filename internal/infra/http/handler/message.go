package handler

import (
	"chatapp/internal/domain/model"
	"chatapp/internal/domain/repository/chatrepo"
	"chatapp/internal/domain/repository/messagerepo"
	"chatapp/internal/domain/repository/userrepo"
	"chatapp/internal/infra/http/request"
	"math/rand"
	"time"

	"github.com/labstack/echo/v4"
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

func (ms *Message) Create(c echo.Context, req request.MessageCreate, uid uint64) (*model.Message, error) {
	id := GenerateMessageID()

	msg := model.Message{
		ID:          id,
		ChatID:      req.ChatID,
		Sender:      uid,
		Receiver:    req.Receiver,
		Content:     req.Content,
		ContentType: req.ContentType,
		CreatedAt:   time.Now(),
	}

	if err := ms.repo.Add(c.Request().Context(), msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
