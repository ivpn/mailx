package model

import (
	"testing"
	"time"
)

func TestSubscriptionActive(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		activeUntil time.Time
		want        bool
	}{
		{
			name:        "future time returns true",
			activeUntil: now.Add(2 * time.Second),
			want:        true,
		},
		{
			name:        "past time returns false",
			activeUntil: now.Add(-2 * time.Second),
			want:        false,
		},
		{
			name:        "equal to now returns false",
			activeUntil: now, // ActiveUntil.After(time.Now()) is strictly greater, so should be false
			want:        false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &Subscription{ActiveUntil: tc.activeUntil}
			got := s.Active()
			if got != tc.want {
				t.Fatalf("Active() = %v, want %v (activeUntil=%v, now=%v)", got, tc.want, tc.activeUntil, time.Now())
			}
		})
	}
}

func TestSubscriptionGracePeriod(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		updatedAt   time.Time
		activeUntil time.Time
		want        bool
	}{
		{
			name:        "outage and within 3-day window => true",
			updatedAt:   now.Add(-49 * time.Hour), // outage (older than 48h)
			activeUntil: now.Add(-48 * time.Hour), // 2 days ago (+3d => +1d > now)
			want:        true,
		},
		{
			name:        "not outage and within 3-day window => false",
			updatedAt:   now.Add(-47 * time.Hour), // not outage
			activeUntil: now.Add(-24 * time.Hour), // 1 day ago (+3d => +2d > now)
			want:        false,
		},
		{
			name:        "outage but outside 3-day window => false",
			updatedAt:   now.Add(-50 * time.Hour),     // outage
			activeUntil: now.Add(-5 * 24 * time.Hour), // 5 days ago (+3d => 2 days before now)
			want:        false,
		},
		{
			name:        "near outage boundary but not outage => false",
			updatedAt:   now.Add(-48 * time.Hour).Add(1 * time.Second), // UpdatedAt +48h ~ now +1s => not outage
			activeUntil: now.Add(-2 * 24 * time.Hour),                  // 2 days ago (+3d => +1d > now)
			want:        false,
		},
		{
			name:        "outage and just inside 3-day window boundary => true",
			updatedAt:   now.Add(-49 * time.Hour),                   // outage
			activeUntil: now.AddDate(0, 0, -3).Add(2 * time.Second), // ActiveUntil +3d ~ now +2s > now
			want:        true,
		},
		{
			name:        "outage and exactly outside 3-day window => false",
			updatedAt:   now.Add(-49 * time.Hour),                    // outage
			activeUntil: now.AddDate(0, 0, -3).Add(-2 * time.Second), // ActiveUntil +3d ~ now -2s < now
			want:        false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &Subscription{
				UpdatedAt:   tc.updatedAt,
				ActiveUntil: tc.activeUntil,
			}
			got := s.GracePeriod()
			if got != tc.want {
				t.Fatalf("GracePeriod() = %v, want %v (updatedAt=%v activeUntil=%v now=%v)", got, tc.want, tc.updatedAt, tc.activeUntil, time.Now())
			}
		})
	}
}

func TestSubscriptionLimitedAccess(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		activeUntil time.Time
		want        bool
	}{
		{
			name:        "active in future => true",
			activeUntil: now.Add(1 * time.Hour),
			want:        true,
		},
		{
			name:        "just expired 1s ago (<14d) => true",
			activeUntil: now.Add(-1 * time.Second),
			want:        true,
		},
		{
			name:        "13d 23h 59m 59s ago (<14d) => true",
			activeUntil: now.Add(-14*24*time.Hour + 1*time.Second),
			want:        true,
		},
		{
			name:        "exactly 14d ago boundary => false",
			activeUntil: now.Add(-14 * 24 * time.Hour),
			want:        false,
		},
		{
			name:        "14d and 1s ago => false",
			activeUntil: now.Add(-14*24*time.Hour - 1*time.Second),
			want:        false,
		},
		{
			name:        "30d ago => false",
			activeUntil: now.Add(-30 * 24 * time.Hour),
			want:        false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &Subscription{ActiveUntil: tc.activeUntil}
			got := s.LimitedAccess()
			if got != tc.want {
				t.Fatalf("LimitedAccess() = %v, want %v (activeUntil=%v now=%v)", got, tc.want, tc.activeUntil, time.Now())
			}
		})
	}
}

func TestSubscriptionPendingDelete(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		activeUntil time.Time
		want        bool
	}{
		{
			name:        "still active in future => false",
			activeUntil: now.Add(1 * time.Hour),
			want:        false,
		},
		{
			name:        "expired 1s ago (<14d) => false",
			activeUntil: now.Add(-1 * time.Second),
			want:        false,
		},
		{
			name:        "13d 23h 59m 59s ago (<14d) => false",
			activeUntil: now.Add(-14*24*time.Hour + 1*time.Second),
			want:        false,
		},
		{
			name:        "exactly 14d ago boundary => true",
			activeUntil: now.Add(-14 * 24 * time.Hour),
			want:        true,
		},
		{
			name:        "14d and 1s ago => true",
			activeUntil: now.Add(-14*24*time.Hour - 1*time.Second),
			want:        true,
		},
		{
			name:        "30d ago => true",
			activeUntil: now.Add(-30 * 24 * time.Hour),
			want:        true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &Subscription{ActiveUntil: tc.activeUntil}
			got := s.PendingDelete()
			if got != tc.want {
				t.Fatalf("PendingDelete() = %v, want %v (activeUntil=%v now=%v)", got, tc.want, tc.activeUntil, time.Now())
			}
		})
	}
}

func TestSubscriptionIsOutage(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		updatedAt time.Time
		want      bool
	}{
		{
			name:      "updated 49h ago => outage",
			updatedAt: now.Add(-49 * time.Hour),
			want:      true,
		},
		{
			name:      "updated 48h + 1s ago => outage",
			updatedAt: now.Add(-48*time.Hour - 1*time.Second),
			want:      true,
		},
		{
			name:      "updated 48h - 1s ago => not outage",
			updatedAt: now.Add(-48*time.Hour + 1*time.Second),
			want:      false,
		},
		{
			name:      "updated 47h ago => not outage",
			updatedAt: now.Add(-47 * time.Hour),
			want:      false,
		},
		{
			name:      "just updated (now) => not outage",
			updatedAt: now,
			want:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &Subscription{UpdatedAt: tc.updatedAt}
			got := s.IsOutage()
			if got != tc.want {
				t.Fatalf("IsOutage() = %v, want %v (updatedAt=%v now=%v threshold=%v)", got, tc.want, tc.updatedAt, time.Now(), tc.updatedAt.Add(48*time.Hour))
			}
		})
	}
}

func TestSubscriptionGracePeriodDays(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		activeUntil time.Time
		days        int
		want        bool
	}{
		{
			name:        "future activeUntil with 3 days window => true",
			activeUntil: now.Add(2 * time.Hour),
			days:        3,
			want:        true,
		},
		{
			name:        "exact boundary 3 days ago => false",
			activeUntil: now.AddDate(0, 0, -3),
			days:        3,
			want:        false,
		},
		{
			name:        "just inside boundary 3 days => true",
			activeUntil: now.AddDate(0, 0, -3).Add(1 * time.Second),
			days:        3,
			want:        true,
		},
		{
			name:        "just outside boundary 3 days => false",
			activeUntil: now.AddDate(0, 0, -3).Add(-1 * time.Second),
			days:        3,
			want:        false,
		},
		{
			name:        "exact boundary 14 days => false",
			activeUntil: now.AddDate(0, 0, -14),
			days:        14,
			want:        false,
		},
		{
			name:        "just inside boundary 14 days => true",
			activeUntil: now.AddDate(0, 0, -14).Add(1 * time.Second),
			days:        14,
			want:        true,
		},
		{
			name:        "past boundary 14 days => false",
			activeUntil: now.AddDate(0, 0, -14).Add(-1 * time.Second),
			days:        14,
			want:        false,
		},
		{
			name:        "exact boundary 1 day => false",
			activeUntil: now.Add(-24 * time.Hour),
			days:        1,
			want:        false,
		},
		{
			name:        "just inside boundary 1 day => true",
			activeUntil: now.Add(-24*time.Hour + 1*time.Second),
			days:        1,
			want:        true,
		},
		{
			name:        "far past even with large window => false",
			activeUntil: now.AddDate(0, 0, -30),
			days:        14,
			want:        false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &Subscription{ActiveUntil: tc.activeUntil}
			got := s.GracePeriodDays(tc.days)
			if got != tc.want {
				t.Fatalf("GracePeriodDays(%d) = %v, want %v (activeUntil=%v now=%v boundary=%v)", tc.days, got, tc.want, tc.activeUntil, time.Now(), tc.activeUntil.AddDate(0, 0, tc.days))
			}
		})
	}
}
