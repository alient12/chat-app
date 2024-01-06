package request

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ContactCreate struct {
	ID   uint64 `json:"id,omitempty" validate:"required"`
	Name string `json:"name,omitempty" validate:"required"`
}

func (cc ContactCreate) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(cc); err != nil {
		return fmt.Errorf("create request validation failed %w", err)
	}

	return nil
}
