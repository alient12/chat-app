package userrepo

import (
	"chatapp/internal/domain/model"
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

var ErrUserIDDuplicate = errors.New("user id already exists")
var ErrUsernameDuplicate = echo.NewHTTPError(http.StatusBadRequest, "username already exists")
var ErrPhoneDuplicate = echo.NewHTTPError(http.StatusBadRequest, "phone number already exists")
var ErrImageSrcDuplicate = errors.New("image source already exists")

type GetCommand struct {
	ID       *uint64
	Username *string
	Phone    *string
	Keyword  *string
}

type Repository interface {
	Add(ctx context.Context, m model.User) error
	Delete(ctx context.Context, id uint64) error
	Update(ctx context.Context, m model.User) error
	Get(ctx context.Context, cmd GetCommand) []model.User
}
