package model

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrDuplicateSubscription = errors.New("subscription already exists")
)

type SubscriptionStatus string

const (
	Active        SubscriptionStatus = "active"
	GracePeriod   SubscriptionStatus = "grace_period"
	LimitedAccess SubscriptionStatus = "limited_access"
	PendingDelete SubscriptionStatus = "pending_delete"
	Tier1         string             = "Tier 1"
)

type Subscription struct {
	ID          string             `gorm:"unique" json:"id"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	UserID      string             `json:"-"`
	ActiveUntil time.Time          `json:"active_until"`
	IsActive    bool               `json:"-"`
	Tier        string             `json:"tier"`
	TokenHash   string             `gorm:"unique" json:"-"`
	Notified    bool               `json:"-"`
	Status      SubscriptionStatus `gorm:"-" json:"status"`
	Outage      bool               `gorm:"-" json:"outage"`
}

func (s *Subscription) Active() bool {
	return s.ActiveUntil.After(time.Now()) && !strings.Contains(s.Tier, Tier1)
}

func (s *Subscription) GracePeriod() bool {
	return s.IsOutage() && s.GracePeriodDays(3)
}

func (s *Subscription) LimitedAccess() bool {
	return s.GracePeriodDays(14)
}

func (s *Subscription) PendingDelete() bool {
	return !s.GracePeriodDays(14)
}

func (s *Subscription) ActiveStatus() bool {
	return s.Active() || s.GracePeriod()
}

func (s *Subscription) IsOutage() bool {
	return s.UpdatedAt.Add(time.Duration(48) * time.Hour).Before(time.Now())
}

func (s *Subscription) GracePeriodDays(days int) bool {
	return s.ActiveUntil.AddDate(0, 0, days).After(time.Now())
}

func (s *Subscription) GetStatus() SubscriptionStatus {
	if s.Active() {
		return Active
	}
	if s.GracePeriod() {
		return GracePeriod
	}
	if s.LimitedAccess() {
		return LimitedAccess
	}
	return PendingDelete
}
