package model

import (
	"errors"
	"regexp"

	"github.com/alexedwards/argon2id"
)

var (
	ErrDuplicateEmail = errors.New("email already exists")
	ErrHashFailed     = errors.New("password hash failed")
	ErrMatchFailed    = errors.New("incorrect password")
	EmailRX           = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type User struct {
	BaseModel
	Email         string  `gorm:"unique" json:"email"`
	PasswordHash  string  `json:"-"`
	PasswordPlain *string `gorm:"-" json:"-"`
	IsActive      bool    `json:"is_active"`
}

type UserStats struct {
	Forwards  int           `json:"forwards"`
	Blocks    int           `json:"blocks"`
	Replies   int           `json:"replies"`
	Sends     int           `json:"sends"`
	Bandwidth int           `json:"bandwidth"`
	Aliases   int64         `json:"aliases"`
	Messages  []interface{} `json:"messages" gorm:"type:text"`
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
