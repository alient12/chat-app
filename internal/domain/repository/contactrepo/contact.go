package contactrepo

import (
	"chatapp/internal/domain/model"
	"context"
	"errors"
)

var ErrContactDuplicate = errors.New("contact already exists")
var ErrContactNotFound = errors.New("contact not found")

type Repository interface {
	Add(ctx context.Context, m model.Contact) error
	Delete(ctx context.Context, uid uint64, cid uint64) error
	Get(ctx context.Context, uid uint64) []model.Contact
}
