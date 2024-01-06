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
	Delete(ctx context.Context, chid uint64, mid uint64) error
	Get(ctx context.Context, cmd GetCommand) []model.Message
}
