package utils

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/alexedwards/argon2id"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmptySecret = errors.New("secret cannot be empty")
	BcryptCost     = bcrypt.DefaultCost
)

func Hash(secret string) (string, error) {
	if strings.TrimSpace(secret) == "" {
		return "", ErrEmptySecret
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(secret), BcryptCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func HashPGPKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

func HashMatches(secret string, hash string) bool {
	if strings.TrimSpace(secret) == "" || strings.TrimSpace(hash) == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(secret))
	return err == nil
}

func HashPassword(secret string) (string, error) {
	if strings.TrimSpace(secret) == "" {
		return "", ErrEmptySecret
	}

	hash, err := argon2id.CreateHash(secret, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func HashMatchesPassword(secret string, hash string) bool {
	if strings.TrimSpace(secret) == "" || strings.TrimSpace(hash) == "" {
		return false
	}

	match, err := argon2id.ComparePasswordAndHash(secret, hash)
	return match && err == nil
}

// TimingSafeEqual performs constant-time comparison to avoid timing attacks
func TimingSafeEqual(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
