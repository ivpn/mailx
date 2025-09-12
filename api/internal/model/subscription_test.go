package model

import (
	"testing"
	"time"
)

func TestSubscription_IsActive(t *testing.T) {
	tests := []struct {
		name        string
		activeUntil time.Time
		want        bool
	}{
		{
			name:        "active subscription",
			activeUntil: time.Now().Add(24 * time.Hour), // 1 day in the future
			want:        true,
		},
		{
			name:        "expired subscription",
			activeUntil: time.Now().Add(-24 * time.Hour), // 1 day in the past
			want:        false,
		},
		{
			name:        "subscription expires now",
			activeUntil: time.Now(),
			want:        false, // time.Now() is not After time.Now()
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Subscription{
				ActiveUntil: tt.activeUntil,
			}
			if got := s.IsActiveCheck(); got != tt.want {
				t.Errorf("Subscription.IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestSubscription_IsActiveWithGracePeriod(t *testing.T) {
	tests := []struct {
		name        string
		activeUntil time.Time
		graceDays   int
		want        bool
	}{
		{
			name:        "active subscription",
			activeUntil: time.Now().Add(24 * time.Hour), // 1 day in the future
			graceDays:   0,
			want:        true,
		},
		{
			name:        "expired subscription but within grace period",
			activeUntil: time.Now().Add(-2 * 24 * time.Hour), // 2 days in the past
			graceDays:   3,
			want:        true,
		},
		{
			name:        "expired subscription outside grace period",
			activeUntil: time.Now().Add(-5 * 24 * time.Hour), // 5 days in the past
			graceDays:   3,
			want:        false,
		},
		{
			name:        "subscription expires today with grace period",
			activeUntil: time.Now(),
			graceDays:   1,
			want:        true,
		},
		{
			name:        "subscription expires today without grace period",
			activeUntil: time.Now(),
			graceDays:   0,
			want:        false, // time.Now() is not After time.Now()
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Subscription{
				ActiveUntil: tt.activeUntil,
			}
			if got := s.IsActiveWithGracePeriod(tt.graceDays); got != tt.want {
				t.Errorf("Subscription.IsActiveWithGracePeriod(%v) = %v, want %v", tt.graceDays, got, tt.want)
			}
		})
	}
}
