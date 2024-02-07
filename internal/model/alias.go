package model

import (
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
)

type Alias struct {
	gorm.Model
	ID          string `json:"id"`
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
