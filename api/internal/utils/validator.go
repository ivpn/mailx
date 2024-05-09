package utils

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	*validator.Validate
}

func NewValidator() Validator {
	return Validator{validator.New()}
}

func ValidateEmail(email string) error {
	validator := NewValidator()
	return validator.Var(email, "required,email")
}
