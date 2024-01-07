package request

import (
	"chatapp/internal/domain/model"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type MessageCreate struct {
	ChatID      uint64                   `json:"chatid,omitempty" validate:"required"`
	Receiver    uint64                   `json:"receiver,omitempty"`
	Content     string                   `json:"content,omitempty" validate:"required"`
	ContentType model.MessageContentType `json:"contenttype,omitempty"`
}

func (msc MessageCreate) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(msc); err != nil {
		return fmt.Errorf("create request validation failed %w", err)
	}

	return nil
}
