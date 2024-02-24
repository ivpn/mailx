package model

import (
	"errors"

	"github.com/alexedwards/argon2id"
)

var (
	ErrHashFailed  = errors.New("password hash failed")
	ErrMatchFailed = errors.New("password match failed")
)

func (u *User) SetPassword(passwordPlain string) error {
	hash, err := argon2id.CreateHash(passwordPlain, argon2id.DefaultParams)
	if err != nil {
		return ErrHashFailed
	}

	u.PasswordHash = hash
	u.PasswordPlain = nil

	return nil
}

func (u *User) Matches(passwordPlain string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(passwordPlain, u.PasswordHash)
	if err != nil {
		return false, ErrMatchFailed
	}

	return match, nil
}
