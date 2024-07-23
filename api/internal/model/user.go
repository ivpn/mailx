package model

import (
	"errors"

	"ivpn.net/email/api/internal/utils"
)

var (
	ErrDuplicateEmail = errors.New("email already exists")
	ErrHashFailed     = errors.New("password hash failed")
)

type User struct {
	BaseModel
	Email          string  `gorm:"unique" json:"email"`
	PasswordHash   string  `json:"-"`
	PasswordPlain  *string `gorm:"-" json:"-"`
	IsActive       bool    `json:"is_active"`
	TotpSecret     string  `json:"-"`
	TotpBackup     string  `json:"-"`
	TotpBackupUsed string  `json:"-"`
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
	hash, err := utils.Hash(passwordPlain)
	if err != nil {
		return ErrHashFailed
	}

	u.PasswordHash = hash
	u.PasswordPlain = nil

	return nil
}

func (u *User) Matches(passwordPlain string) bool {
	return utils.HashMatches(passwordPlain, u.PasswordHash)
}

func (u *User) IsTotpEnabled() bool {
	return u.TotpSecret != ""
}
