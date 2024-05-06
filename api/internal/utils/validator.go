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
	if email == "" || !EmailRX.MatchString(email) {
		return ErrInvalidEmail
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" || len(password) < 8 || len(password) > 64 {
		return ErrInvalidPassword
	}

	return nil
}
