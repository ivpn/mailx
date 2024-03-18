package model

import "time"

type SubscriptionType string

const (
	Free    SubscriptionType = "Free"
	Managed SubscriptionType = "Managed"
)

type Subscription struct {
	BaseModel
	UserID       string           `json:"-"`
	ActivationID string           `json:"activation_id"`
	Type         SubscriptionType `json:"type"`
	ActiveUntil  *time.Time       `json:"active_until"`
}
