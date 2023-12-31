package chatrepo

import (
	"chatapp/internal/domain/model"
	"context"
)

type GetCommand struct {
	ID *uint64
}

type Repository interface {
	Add(ctx context.Context, m model.Chat) error
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, cmd GetCommand) []model.Chat
}
