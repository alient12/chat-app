package filerepo

import (
	"chatapp/internal/domain/model"
	"context"
)

type GetCommand struct {
	ID          *uint64
	UserID      *uint64
	FileName    *uint64
	ContentType *string
	ChatID      *uint64
}

type Repository interface {
	Add(ctx context.Context, m model.File) error
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, cmd GetCommand) []model.File
	Update(ctx context.Context, m model.File) error
}
