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
	UserID     string    `json:"user_id"`
	TokenHash  string    `json:"-"`
	TokenPlain *string   `gorm:"-" json:"-"`
	Name       string    `json:"name"`
	ExpiresAt  time.Time `json:"expires_at"`
}

func (u *AccessKey) SetToken(tokenPlain string) error {
	hash, err := utils.HashPassword(tokenPlain)
	if err != nil {
		return ErrTokenHashFailed
	}

	u.TokenHash = hash
	u.TokenPlain = nil

	return nil
}

func (u *AccessKey) Matches(tokenPlain string) bool {
	return utils.HashMatchesPassword(tokenPlain, u.TokenHash)
}
