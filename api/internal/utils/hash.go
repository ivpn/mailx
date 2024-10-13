package utils

import (
	"github.com/alexedwards/argon2id"
	"golang.org/x/crypto/bcrypt"
)

func Hash(secret string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func HashMatches(secret string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(secret))
	return err == nil
}

func HashPassword(secret string) (string, error) {
	hash, err := argon2id.CreateHash(secret, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func HashMatchesPassword(secret string, hash string) bool {
	match, err := argon2id.ComparePasswordAndHash(secret, hash)
	return match && err == nil
}
