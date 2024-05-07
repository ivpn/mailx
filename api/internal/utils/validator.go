package utils

import (
	"errors"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidOTP      = errors.New("invalid otp")
	EmailRX            = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
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
