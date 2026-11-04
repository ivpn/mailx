package model

import (
	"testing"
	"time"
)

func TestAliasIsPaused(t *testing.T) {
	now := time.Date(2026, 10, 2, 12, 0, 0, 0, time.UTC)
	at := func(d time.Duration) *time.Time {
		v := now.Add(d)
		return &v
	}

	tests := []struct {
		name        string
		pausedUntil *time.Time
		want        bool
	}{
		{name: "never paused", pausedUntil: nil, want: false},
		{name: "pause ended", pausedUntil: at(-time.Second), want: false},
		{name: "pause ends exactly now", pausedUntil: at(0), want: false},
		{name: "paused", pausedUntil: at(time.Second), want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alias := Alias{PausedUntil: tt.pausedUntil}
			if got := alias.IsPaused(now); got != tt.want {
				t.Errorf("IsPaused() = %v, want %v", got, tt.want)
			}

			alias.ClearExpiredPause(now)
			if tt.want && alias.PausedUntil != tt.pausedUntil {
				t.Errorf("ClearExpiredPause() changed a running pause to %v", alias.PausedUntil)
			}
			if !tt.want && alias.PausedUntil != nil {
				t.Errorf("ClearExpiredPause() kept %v, want nil", alias.PausedUntil)
			}
		})
	}
}

func TestPauseDuration(t *testing.T) {
	tests := []struct {
		value int
		unit  string
		want  time.Duration
	}{
		{value: 1, unit: PauseHours, want: time.Hour},
		{value: 10, unit: PauseHours, want: 10 * time.Hour},
		{value: 1, unit: PauseDays, want: 24 * time.Hour},
		{value: 10, unit: PauseDays, want: 240 * time.Hour},
		{value: 1, unit: PauseWeeks, want: 7 * 24 * time.Hour},
		{value: 10, unit: PauseWeeks, want: 70 * 24 * time.Hour},
	}

	for _, tt := range tests {
		got, err := PauseDuration(tt.value, tt.unit)
		if err != nil || got != tt.want {
			t.Errorf("PauseDuration(%d, %q) = %v, %v, want %v", tt.value, tt.unit, got, err, tt.want)
		}
	}

	invalid := []struct {
		value int
		unit  string
	}{
		{value: 0, unit: PauseDays},
		{value: 11, unit: PauseDays},
		{value: -1, unit: PauseHours},
		{value: 1, unit: "months"},
		{value: 1, unit: ""},
		{value: 1, unit: "Days"},
	}

	for _, tt := range invalid {
		if _, err := PauseDuration(tt.value, tt.unit); err != ErrInvalidPauseDuration {
			t.Errorf("PauseDuration(%d, %q) error = %v, want ErrInvalidPauseDuration", tt.value, tt.unit, err)
		}
	}
}
