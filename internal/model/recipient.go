package model

import (
	"errors"
	"math/rand"
	"time"
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

func GenerateVerification() string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	const charset = "0123456789"
	result := make([]byte, 6)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}

	return string(result)
}
