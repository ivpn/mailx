package model

import (
	"errors"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type User struct {
	UUIDBaseModel
	Email    string   `json:"email"`
	Password password `json:"-"`
}

type password struct {
	plaintext *string
	hash      string
}
