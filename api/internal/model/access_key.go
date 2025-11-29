package model

import (
	"crypto/rand"
	"errors"
	"time"

	"ivpn.net/email/api/internal/utils"
)

var (
	ErrTokenHashFailed = errors.New("token hash failed")
)

type AccessKey struct {
	BaseModel
	UserId     string     `json:"user_id"`
	TokenHash  string     `json:"-"`
	TokenId    string     `gorm:"-" json:"-"`
	TokenPlain *string    `gorm:"-" json:"-"`
	Name       string     `json:"name"`
	ExpiresAt  *time.Time `json:"expires_at"` // nullable
}

func (a *AccessKey) SetToken(token string) error {
	hash, err := utils.HashPassword(token)
	if err != nil {
		return ErrTokenHashFailed
	}

	a.TokenHash = hash
	a.TokenPlain = nil

	return nil
}

func (a *AccessKey) Matches(tokenPlain string) bool {
	return utils.HashMatchesPassword(tokenPlain, a.TokenHash)
}

func NeverExpires() *time.Time {
	return nil
}

func (a *AccessKey) IsExpired() bool {
	if a.ExpiresAt == nil {
		return false // never expires
	}

	return time.Now().After(*a.ExpiresAt)
}

func GenAccessKeyToken() (string, error) {
	token, err := GenToken(48)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GenToken(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}

	return string(b), nil
}
