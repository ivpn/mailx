package model

import "time"

type Subscription struct {
	BaseModel
	Name         string     `json:"name"`
	UserID       string     `json:"-"`
	ActivationID string     `json:"activation_id"`
	IsActive     bool       `json:"is_active"`
	ActiveUntil  *time.Time `json:"active_until"`
}
