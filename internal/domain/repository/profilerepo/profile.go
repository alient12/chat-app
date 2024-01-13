package profilerepo

import (
	"chatapp/internal/domain/model"
	"context"
	"errors"
)

var ErrIDDuplicate = errors.New("duplicated user id")
var ErrIDNotFound = errors.New("user id not found")

type GetCommand struct {
	UserID *uint64
}

type Repository interface {
	Add(ctx context.Context, m model.Profile) error
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, cmd GetCommand) []model.Profile
}
