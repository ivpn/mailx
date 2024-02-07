package model

import "github.com/jinzhu/gorm"

type Recipient struct {
	gorm.Model
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	Verification string `json:"verification"`
	Verified     bool   `json:"verified"`
}
