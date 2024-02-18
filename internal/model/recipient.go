package model

import (
	"math/rand"
	"time"
)

type Recipient struct {
	BaseModel
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	Verification string `json:"verification"`
	Verified     bool   `json:"verified"`
}

func GenerateVerification() string {
	rand.Seed(time.Now().UnixNano())

	const charset = "0123456789"
	result := make([]byte, 6)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}
