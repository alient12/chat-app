package request

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ChatCreate struct {
	People []uint64 `json:"people,omitempty" validate:"required"`
	Token  string   `json:"token,omitempty"`
}

func (chc ChatCreate) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(chc); err != nil {
		return fmt.Errorf("create request validation failed %w", err)
	}

	return nil
}
