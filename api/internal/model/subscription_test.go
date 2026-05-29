package model

import (
	"testing"
	"time"
)

// helpers to build a Subscription at a point in time
func activeSubscription() *Subscription {
	return &Subscription{
		ActiveUntil: time.Now().Add(30 * 24 * time.Hour), // 30 days from now
		UpdatedAt:   time.Now().Add(-1 * time.Hour),      // updated 1 hour ago (no outage)
		Tier:        "Tier 2",
	}
}

func expiredSubscription() *Subscription {
	return &Subscription{
		ActiveUntil: time.Now().Add(-30 * 24 * time.Hour), // expired 30 days ago
		UpdatedAt:   time.Now().Add(-1 * time.Hour),
		Tier:        "Tier 2",
	}
}

// --- Active ---

func TestActive(t *testing.T) {
	future := time.Now().Add(30 * 24 * time.Hour)
	recent := time.Now().Add(-1 * time.Hour)
	outageTime := time.Now().Add(-49 * time.Hour)

	tests := []struct {
		name        string
		activeUntil time.Time
		updatedAt   time.Time
		tier        string
		terminated  bool
		want        bool
	}{
		{
			name:        "active: future ActiveUntil, non-Tier1, no outage",
			activeUntil: future,
			updatedAt:   recent,
			tier:        "Tier 2",
			want:        true,
		},
		{
			name:        "not active: ActiveUntil in the past",
			activeUntil: time.Now().Add(-1 * time.Hour),
			updatedAt:   recent,
			tier:        "Tier 2",
			want:        false,
		},
		{
			name:        "not active: tier is IVPN Tier 1",
			activeUntil: future,
			updatedAt:   recent,
			tier:        Tier1,
			want:        false,
		},
		{
			name:        "not active: outage (UpdatedAt > 48h ago)",
			activeUntil: future,
			updatedAt:   outageTime,
			tier:        "Tier 2",
			want:        false,
		},
		{
			name:        "not active: expired and outage and Tier1",
			activeUntil: time.Now().Add(-1 * time.Hour),
			updatedAt:   outageTime,
			tier:        Tier1,
			want:        false,
		},
		{
			name:        "not active: UpdatedAt is zero (no outage), but ActiveUntil in past",
			activeUntil: time.Now().Add(-1 * time.Hour),
			updatedAt:   time.Time{},
			tier:        "Tier 2",
			want:        false,
		},
		{
			name:        "active: UpdatedAt is zero (no outage), ActiveUntil in future, non-Tier1",
			activeUntil: future,
			updatedAt:   time.Time{},
			tier:        "Tier 2",
			want:        true,
		},
		{
			name:        "not active: terminated, even with otherwise valid subscription",
			activeUntil: future,
			updatedAt:   recent,
			tier:        "Tier 2",
			terminated:  true,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Subscription{
				ActiveUntil: tt.activeUntil,
				UpdatedAt:   tt.updatedAt,
				Tier:        tt.tier,
				Terminated:  tt.terminated,
			}
			if got := s.Active(); got != tt.want {
				t.Errorf("Active() = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- GracePeriodDays ---

func TestGracePeriodDays(t *testing.T) {
	tests := []struct {
		name        string
		activeUntil time.Time
		days        int
		want        bool
	}{
		{
			name:        "active until yesterday + 3 days grace = still in grace",
			activeUntil: time.Now().Add(-24 * time.Hour),
			days:        3,
			want:        true,
		},
		{
			name:        "active until 4 days ago + 3 days grace = outside grace",
			activeUntil: time.Now().AddDate(0, 0, -4),
			days:        3,
			want:        false,
		},
		{
			name:        "active until tomorrow + 0 days = not in grace (not yet expired)",
			activeUntil: time.Now().Add(24 * time.Hour),
			days:        0,
			want:        false,
		},
		{
			name:        "active until 15 days ago + 14 days grace = outside grace",
			activeUntil: time.Now().AddDate(0, 0, -15),
			days:        14,
			want:        false,
		},
		{
			name:        "active until 13 days ago + 14 days grace = in grace",
			activeUntil: time.Now().AddDate(0, 0, -13),
			days:        14,
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Subscription{ActiveUntil: tt.activeUntil}
			if got := s.GracePeriodDays(tt.days); got != tt.want {
				t.Errorf("GracePeriodDays(%d) = %v, want %v", tt.days, got, tt.want)
			}
		})
	}
}

// --- OutageGracePeriodDays ---

func TestOutageGracePeriodDays(t *testing.T) {
	tests := []struct {
		name      string
		updatedAt time.Time
		days      int
		want      bool
	}{
		{
			name:      "updated 1 day ago + 3 days = in outage grace",
			updatedAt: time.Now().Add(-24 * time.Hour),
			days:      3,
			want:      true,
		},
		{
			name:      "updated 4 days ago + 3 days = outside outage grace",
			updatedAt: time.Now().AddDate(0, 0, -4),
			days:      3,
			want:      false,
		},
		{
			name:      "updated 13 days ago + 14 days = in outage grace",
			updatedAt: time.Now().AddDate(0, 0, -13),
			days:      14,
			want:      true,
		},
		{
			name:      "updated 15 days ago + 14 days = outside outage grace",
			updatedAt: time.Now().AddDate(0, 0, -15),
			days:      14,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Subscription{UpdatedAt: tt.updatedAt}
			if got := s.OutageGracePeriodDays(tt.days); got != tt.want {
				t.Errorf("OutageGracePeriodDays(%d) = %v, want %v", tt.days, got, tt.want)
			}
		})
	}
}

// --- GracePeriod ---

func TestGracePeriod(t *testing.T) {
	future := time.Now().Add(30 * 24 * time.Hour)
	recent := time.Now().Add(-1 * time.Hour)
	outageTime := time.Now().Add(-49 * time.Hour)
	outageOutside3Days := time.Now().AddDate(0, 0, -4) // outage AND outside 3-day grace

	tests := []struct {
		name        string
		activeUntil time.Time
		updatedAt   time.Time
		terminated  bool
		want        bool
	}{
		{
			name:        "grace period: outage + within 3-day ActiveUntil grace + within 3-day outage grace",
			activeUntil: time.Now().Add(-24 * time.Hour), // expired 1d ago, within 3d grace
			updatedAt:   outageTime,                      // outage, but < 3d
			want:        true,
		},
		{
			name:        "no grace period: no outage",
			activeUntil: future,
			updatedAt:   recent,
			want:        false,
		},
		{
			name:        "no grace period: outage but ActiveUntil grace expired",
			activeUntil: time.Now().AddDate(0, 0, -4), // 4 days ago, outside 3d grace
			updatedAt:   outageTime,
			want:        false,
		},
		{
			name:        "no grace period: outage but OutageGracePeriodDays(3) expired",
			activeUntil: time.Now().Add(-24 * time.Hour),
			updatedAt:   outageOutside3Days,
			want:        false,
		},
		{
			name:        "no grace period: terminated, even with otherwise valid grace period conditions",
			activeUntil: time.Now().Add(-24 * time.Hour), // expired 1d ago, within 3d grace
			updatedAt:   outageTime,                      // outage, but < 3d
			terminated:  true,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Subscription{
				ActiveUntil: tt.activeUntil,
				UpdatedAt:   tt.updatedAt,
				Terminated:  tt.terminated,
			}
			if got := s.GracePeriod(); got != tt.want {
				t.Errorf("GracePeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- ActiveStatus ---

func TestActiveStatus(t *testing.T) {
	future := time.Now().Add(30 * 24 * time.Hour)
	recent := time.Now().Add(-1 * time.Hour)
	outageTime := time.Now().Add(-49 * time.Hour)

	tests := []struct {
		name        string
		activeUntil time.Time
		updatedAt   time.Time
		tier        string
		want        bool
	}{
		{
			name:        "active status: subscription is active",
			activeUntil: future,
			updatedAt:   recent,
			tier:        "Tier 2",
			want:        true,
		},
		{
			name:        "active status: subscription is in grace period",
			activeUntil: time.Now().Add(-24 * time.Hour), // expired 1d ago, within 3d grace
			updatedAt:   outageTime,
			tier:        "Tier 2",
			want:        true,
		},
		{
			name:        "not active: limited access (no outage, expired)",
			activeUntil: time.Now().AddDate(0, 0, -5),
			updatedAt:   recent,
			tier:        "Tier 2",
			want:        false,
		},
		{
			name:        "not active: pending delete",
			activeUntil: time.Now().AddDate(0, 0, -15),
			updatedAt:   time.Now().AddDate(0, 0, -15),
			tier:        "Tier 2",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Subscription{
				ActiveUntil: tt.activeUntil,
				UpdatedAt:   tt.updatedAt,
				Tier:        tt.tier,
			}
			if got := s.ActiveStatus(); got != tt.want {
				t.Errorf("ActiveStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsOutage(t *testing.T) {
	tests := []struct {
		name      string
		updatedAt time.Time
		want      bool
	}{
		{
			name:      "updated more than 48h ago is outage",
			updatedAt: time.Now().Add(-49 * time.Hour),
			want:      true,
		},
		{
			name:      "updated exactly 48h ago is outage",
			updatedAt: time.Now().Add(-48 * time.Hour).Add(-time.Second),
			want:      true,
		},
		{
			name:      "updated less than 48h ago is not outage",
			updatedAt: time.Now().Add(-1 * time.Hour),
			want:      false,
		},
		{
			name:      "updated just now is not outage",
			updatedAt: time.Now(),
			want:      false,
		},
		{
			name:      "zero time is not outage",
			updatedAt: time.Time{},
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Subscription{UpdatedAt: tt.updatedAt}
			if got := s.IsOutage(); got != tt.want {
				t.Errorf("IsOutage() = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- LimitedAccess ---

func TestLimitedAccess(t *testing.T) {
	tests := []struct {
		name        string
		updatedAt   time.Time
		activeUntil time.Time
		want        bool
	}{
		{
			name:        "limited: updatedAt more than 3 days ago",
			updatedAt:   time.Now().AddDate(0, 0, -4),
			activeUntil: time.Now().Add(30 * 24 * time.Hour),
			want:        true,
		},
		{
			name:        "limited: activeUntil more than 3 days ago",
			updatedAt:   time.Now().Add(-1 * time.Hour),
			activeUntil: time.Now().AddDate(0, 0, -4),
			want:        true,
		},
		{
			name:        "not limited: updatedAt 2 days ago, activeUntil in future",
			updatedAt:   time.Now().AddDate(0, 0, -2),
			activeUntil: time.Now().Add(30 * 24 * time.Hour),
			want:        false,
		},
		{
			name:        "not limited: both updatedAt and activeUntil within 3-day window",
			updatedAt:   time.Now().AddDate(0, 0, -2),
			activeUntil: time.Now().AddDate(0, 0, -2),
			want:        false,
		},
		{
			name:        "not limited: updatedAt recent, activeUntil in future",
			updatedAt:   time.Now().Add(-1 * time.Hour),
			activeUntil: time.Now().Add(30 * 24 * time.Hour),
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Subscription{
				UpdatedAt:   tt.updatedAt,
				ActiveUntil: tt.activeUntil,
			}
			if got := s.LimitedAccess(); got != tt.want {
				t.Errorf("LimitedAccess() = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- GetStatus ---

func TestGetStatus(t *testing.T) {
	future := time.Now().Add(30 * 24 * time.Hour)
	recent := time.Now().Add(-1 * time.Hour)
	outageTime := time.Now().Add(-49 * time.Hour)

	tests := []struct {
		name        string
		activeUntil time.Time
		updatedAt   time.Time
		tier        string
		want        SubscriptionStatus
	}{
		{
			name:        "active subscription returns Active",
			activeUntil: future,
			updatedAt:   recent,
			tier:        "Tier 2",
			want:        Active,
		},
		{
			name:        "grace period subscription returns GracePeriod",
			activeUntil: time.Now().Add(-24 * time.Hour), // expired 1d ago, within 3d grace
			updatedAt:   outageTime,                      // outage, within 3d outage grace
			tier:        "Tier 2",
			want:        GracePeriod,
		},
		{
			name:        "expired subscription with no outage returns LimitedAccess",
			activeUntil: time.Now().AddDate(0, 0, -5),
			updatedAt:   recent,
			tier:        "Tier 2",
			want:        LimitedAccess,
		},
		{
			name:        "IVPN Tier 1 subscription (not active) returns LimitedAccess",
			activeUntil: future,
			updatedAt:   recent,
			tier:        Tier1,
			want:        LimitedAccess,
		},
		{
			name:        "expired subscription outside grace returns LimitedAccess",
			activeUntil: time.Now().AddDate(0, 0, -4), // outside 3d grace
			updatedAt:   outageTime,                   // outage but ActiveUntil grace expired
			tier:        "Tier 2",
			want:        LimitedAccess,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Subscription{
				ActiveUntil: tt.activeUntil,
				UpdatedAt:   tt.updatedAt,
				Tier:        tt.tier,
			}
			if got := s.GetStatus(); got != tt.want {
				t.Errorf("GetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
