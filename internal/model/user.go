package model

import (
	"errors"
	"regexp"

	"github.com/alexedwards/argon2id"
)

var (
	ErrEmailInvalid    = errors.New("invalid email")
	ErrPasswordInvalid = errors.New("invalid password")
	ErrDuplicateEmail  = errors.New("duplicate email")
	ErrHashFailed      = errors.New("password hash failed")
	ErrMatchFailed     = errors.New("password match failed")
	EmailRX            = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type User struct {
	BaseModel
	Email         string  `gorm:"unique" json:"email"`
	PasswordHash  string  `json:"-"`
	PasswordPlain *string `gorm:"-" json:"-"`
	IsActive      bool    `json:"is_active"`
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

func (u *User) SetPassword(passwordPlain string) error {
	hash, err := argon2id.CreateHash(passwordPlain, argon2id.DefaultParams)
	if err != nil {
		return ErrHashFailed
	}

	u.PasswordHash = hash
	u.PasswordPlain = nil

	return nil
}

func (u *User) Matches(passwordPlain string) bool {
	match, err := argon2id.ComparePasswordAndHash(passwordPlain, u.PasswordHash)
	if err != nil {
		return false
	}

	return match
}
