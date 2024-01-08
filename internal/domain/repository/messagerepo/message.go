package messagerepo

import (
	"chatapp/internal/domain/model"
	"context"
)

type GetCommand struct {
	ID          *uint64
	ChatID      *uint64
	Sender      *uint64
	Keyword     *string
	ContentType *model.MessageContentType
}

type Repository interface {
	Add(ctx context.Context, m model.Message) error
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, cmd GetCommand) []model.Message
}
