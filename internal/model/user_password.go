package model

import (
	"log"

	"github.com/alexedwards/argon2id"
)

func (u *User) Set(plaintextPassword string) error {
	hash, err := argon2id.CreateHash(plaintextPassword, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	u.PasswordPlain = &plaintextPassword
	u.PasswordHash = hash

	return nil
}

func (u *User) Matches(plaintextPassword string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(plaintextPassword, u.PasswordHash)
	if err != nil {
		log.Fatal(err)
	}

	return match, nil
}
