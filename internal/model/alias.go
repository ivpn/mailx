package model

import (
	"math/rand"
	"time"
)

type Alias struct {
	PK          uint   `gorm:"primaryKey"`
	ID          string `json:"id"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	RecipientID string `json:"recipient_id"`
	Slug        string `json:"slug"`
}

func GenerateSlug() string {
	rand.Seed(time.Now().UnixNano())

	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 8)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}
