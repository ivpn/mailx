package model

import (
	"math/rand"
	"time"
)

type Alias struct {
	UUIDBaseModel
	RecipientID string `json:"recipient_id"`
	Name        string `json:"name"`
}

func GenerateName() string {
	rand.Seed(time.Now().UnixNano())

	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 8)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}
