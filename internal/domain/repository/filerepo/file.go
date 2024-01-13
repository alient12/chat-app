package filerepo

import (
	"chatapp/internal/domain/model"
	"context"
	"errors"
)

var ErrFileIDDuplicate = errors.New("duplicated file id")
var ErrIDNotFound = errors.New("file id not found")

type GetCommand struct {
	ID          *uint64
	UserID      *uint64
	FileName    *string
	ContentType *string
	ChatID      *uint64
	Keyword     *string
}

type Repository interface {
	Add(ctx context.Context, m model.File) error
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, cmd GetCommand) []model.File
	Update(ctx context.Context, m model.File) error
}
