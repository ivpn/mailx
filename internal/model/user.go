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
	BaseModel
	Email         string  `json:"email"`
	PasswordHash  string  `json:"-"`
	PasswordPlain *string `gorm:"-"`
}

func (u *User) Validate() error {
	if u.Email == "" || !EmailRX.MatchString(u.Email) {
		return ErrEmailInvalid
	}

	if u.PasswordPlain == nil || len(*u.PasswordPlain) < 8 || len(*u.PasswordPlain) > 64 {
		return ErrPasswordInvalid
	}

	return nil
}
