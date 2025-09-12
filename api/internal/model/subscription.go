package model

import (
	"errors"
	"time"
)

var (
	ErrDuplicateSubscription = errors.New("subscription already exists")
)

type Subscription struct {
	ID          string    `gorm:"unique" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"-"`
	UserID      string    `json:"-"`
	ActiveUntil time.Time `json:"active_until"`
	IsActive    bool      `json:"is_active"`
	Tier        string    `json:"tier"`
	TokenHash   string    `gorm:"unique" json:"-"`
}

func (s *Subscription) IsActiveCheck() bool {
	return s.ActiveUntil.After(time.Now()) || s.IsActive
}

func (s *Subscription) IsActiveWithGracePeriod(days int) bool {
	return s.ActiveUntil.AddDate(0, 0, days).After(time.Now())
}
