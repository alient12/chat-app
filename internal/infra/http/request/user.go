package request

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type UserCreate struct {
	Username  string `json:"username,omitempty" validate:"required"`
	Password  string `json:"password,omitempty" validate:"required,min=8"`
	Phone     string `json:"phone,omitempty" validate:"required,len=10,numeric"`
	Firstname string `json:"firstname,omitempty" validate:"required,alpha"`
	Lastname  string `json:"lastname,omitempty" validate:"required,alpha"`
	Image     string `json:"image,omitempty"`
	Bio       string `json:"bio,omitempty" validate:"max=100"`
	Token     string `json:"token,omitempty"`
}

func (uc UserCreate) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(uc); err != nil {
		return fmt.Errorf("create request validation failed %w", err)
	}

	return nil
}

type UserLogin struct {
	Username string `json:"username,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required,min=8"`
}

func (uc UserLogin) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(uc); err != nil {
		return fmt.Errorf("create request validation failed %w", err)
	}

	return nil
}
