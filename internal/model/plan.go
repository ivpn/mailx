package model

import "time"

type Plan struct {
	BaseModel
	Name         string     `json:"name"`
	UserID       string     `json:"user_id"`
	ActivationID string     `json:"id"`
	Active       bool       `json:"active"`
	ActiveUntil  *time.Time `json:"active_until"`
}
