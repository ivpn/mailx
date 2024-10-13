package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	*validator.Validate
}

func NewValidator() Validator {
	v := Validator{validator.New()}
	v.RegisterValidation("password", passwordValidation)
	return v
}

func ValidateEmail(email string) error {
	validator := NewValidator()
	return validator.Var(email, "required,email")
}

func passwordValidation(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 12 || len(password) > 64 {
		return false
	}

	var uppercase = regexp.MustCompile(`[A-Z]`).MatchString
	var lowercase = regexp.MustCompile(`[a-z]`).MatchString
	var number = regexp.MustCompile(`[0-9]`).MatchString
	var specialChar = regexp.MustCompile(`[!@#$%^&*(),;.?":{}|<>]`).MatchString

	if !uppercase(password) {
		return false
	}

	if !lowercase(password) {
		return false
	}

	if !number(password) {
		return false
	}

	if !specialChar(password) {
		return false
	}

	return true
}
