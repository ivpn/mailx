package model

import "time"

type Preauth struct {
	ID          string    `json:"id"`
	TokenHash   string    `json:"token_hash"`
	IsActive    bool      `json:"is_active"`
	ActiveUntil time.Time `json:"active_until"`
	Tier        string    `json:"tier"`
}
