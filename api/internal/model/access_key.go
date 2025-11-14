package model

import (
	"errors"
	"time"

	"ivpn.net/email/api/internal/utils"
)

var (
	ErrTokenHashFailed = errors.New("token hash failed")
)

type AccessKey struct {
	BaseModel
	UserID     string     `json:"user_id"`
	TokenHash  string     `json:"-"`
	TokenPlain *string    `gorm:"-" json:"-"`
	Name       string     `json:"name"`
	ExpiresAt  *time.Time `json:"expires_at"` // nullable
}

func (a *AccessKey) SetToken(tokenPlain string) error {
	hash, err := utils.HashPassword(tokenPlain)
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
