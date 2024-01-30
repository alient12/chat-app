package chatrepo

import (
	"chatapp/internal/domain/model"
	"context"
	"errors"
)

var ErrChatIDDuplicate = errors.New("duplicated chat id")
var ErrIDNotFound = errors.New("user id not found")
var ErrDualChatDuplicate = errors.New("duplicated dual chat")

type GetCommand struct {
	ID     *uint64
	UserID *[]uint64
}

type Repository interface {
	Add(ctx context.Context, m model.Chat) error
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, cmd GetCommand) []model.Chat
	Update(ctx context.Context, m model.Chat) error
}
