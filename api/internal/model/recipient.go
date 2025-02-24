package model

import (
	"errors"
	"strings"
)

var (
	ErrDuplicateRecipient = errors.New("email already exists")
)

type Recipient struct {
	BaseModel
	UserID     string `json:"-"`
	Email      string `gorm:"unique" json:"email"`
	IsActive   bool   `json:"is_active"`
	PGPKey     string `json:"pgp_key"`
	PGPEnabled bool   `json:"pgp_enabled"`
	PGPInline  bool   `json:"pgp_inline"`
}

func GetEmails(rcps []Recipient) string {
	emails := []string{}
	for _, r := range rcps {
		emails = append(emails, r.Email)
	}

	return strings.Join(emails, ",")
}
