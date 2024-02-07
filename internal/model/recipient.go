package model

import "time"

type Recipient struct {
	ID           string `json:"id"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	Verification string `json:"verification"`
	Verified     bool   `json:"verified"`
}
