package model

import (
	"errors"

	"ivpn.net/email-service/internal/utils"
)

var (
	ErrDuplicateRecipient = errors.New("duplicate recipient")
)

type Recipient struct {
	BaseModel
	UserID       string `json:"user_id"`
	Email        string `gorm:"unique" json:"email"`
	Verification string `json:"verification"`
	Verified     bool   `json:"verified"`
}

func (r *Recipient) Validate() error {
	err := utils.ValidateEmail(r.Email)
	if err != nil {
		return err
	}

	return nil
}
