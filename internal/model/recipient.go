package model

import (
	"errors"

	"ivpn.net/email-service/internal/utils"
)

var (
	ErrDuplicateRecipient = errors.New("email already exists")
)

type Recipient struct {
	BaseModel
	UserID   string `json:"-"`
	Email    string `gorm:"unique" json:"email"`
	IsActive bool   `json:"is_active"`
}

func (r *Recipient) Validate() error {
	err := utils.ValidateEmail(r.Email)
	if err != nil {
		return err
	}

	return nil
}
