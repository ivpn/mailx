package model

import (
	"errors"
	"regexp"
)

var (
	ErrEmailInvalid    = errors.New("invalid email")
	ErrPasswordInvalid = errors.New("invalid password")
	EmailRX            = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
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

func (u *User) IsValid() error {
	if u.Email == "" || !EmailRX.MatchString(u.Email) {
		return ErrEmailInvalid
	}

	if u.Password.plaintext == nil || len(*u.Password.plaintext) < 8 || len(*u.Password.plaintext) > 64 {
		return ErrPasswordInvalid
	}

	return nil
}
