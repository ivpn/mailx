package model

import (
	"errors"
	"time"
)

type SubscriptionType string

const (
	Free    SubscriptionType = "Free"
	Managed SubscriptionType = "Managed"
)

var (
	ErrDuplicateSubscription = errors.New("subscription already exists")
)

type Subscription struct {
	ID          string           `gorm:"unique" json:"id"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"-"`
	UserID      string           `json:"-"`
	Type        SubscriptionType `json:"type"`
	ActiveUntil time.Time        `json:"active_until"`
}
