package handler

import (
	"chatapp/internal/domain/model"
	"chatapp/internal/domain/repository/chatrepo"
	"chatapp/internal/domain/repository/messagerepo"
	"chatapp/internal/domain/repository/userrepo"
	"chatapp/internal/infra/http/request"
	"chatapp/internal/util"
	"math/rand"
	"net/http"
	"strconv"
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

func (ms *Message) Get(c echo.Context) error {
	var chatIDPtr, idPtr *uint64
	if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
		chatIDPtr = &id
	} else {
		return echo.ErrBadRequest
	}

	// check auth
	if ckID, _, err := CheckJWT(c); err != nil {
		return err
	} else {
		idPtr = &ckID
	}

	chat := ms.chrepo.Get(c.Request().Context(), chatrepo.GetCommand{
		ID:     chatIDPtr,
		UserID: nil,
	})[0]

	// check if user has access to the chat
	if !util.InSlice(chat.People, *idPtr) {
		return echo.ErrForbidden
	}

	messages := ms.repo.Get(c.Request().Context(), messagerepo.GetCommand{
		ID:          nil,
		ChatID:      chatIDPtr,
		Sender:      nil,
		Keyword:     nil,
		ContentType: nil,
	})

	chat_messages := struct {
		Chat     model.Chat
		Messages []model.Message
	}{
		Chat:     chat,
		Messages: messages,
	}

	return c.JSON(http.StatusOK, chat_messages)
}

func (ms *Message) Delete(c echo.Context) error {
	var chatIDPtr, idPtr *uint64

	if id, err := strconv.ParseUint(c.Param("chid"), 10, 64); err == nil {
		chatIDPtr = &id
	} else {
		return echo.ErrBadRequest
	}

	if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
		idPtr = &id
	} else {
		return echo.ErrBadRequest
	}

	// check auth
	if ckID, _, err := CheckJWT(c); err != nil {
		return err
	} else {
		if ckID != *idPtr {
			return echo.ErrUnauthorized
		}
	}

	if err := ms.repo.Delete(c.Request().Context(), *idPtr); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, *chatIDPtr)
}

func (ms *Message) Register(g *echo.Group) {
	g.GET("/chats/:id", ms.Get)
	g.DELETE("/chats/:chid/messages/:id", ms.Delete)
}
