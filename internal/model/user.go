package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type User struct {
	PK        uint      `gorm:"primaryKey"`
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string   `json:"email"`
	Password  password `json:"-"`
}

type password struct {
	plaintext *string
	hash      string
}
