package repository

import (
	"context"
	"testing"
	"time"

	"ivpn.net/email/api/internal/model"
)

func createTestAlias(t *testing.T, d *Database, name string, userID string) model.Alias {
	t.Helper()

	alias := model.Alias{Name: name, UserID: userID, Enabled: true, Description: "Shopping", Recipients: "me@example.com"}
	if err := d.Client.Create(&alias).Error; err != nil {
		t.Fatalf("create alias %s: %v", name, err)
	}

	return alias
}

func storedPausedUntil(t *testing.T, d *Database, ID string) *time.Time {
	t.Helper()

	var alias model.Alias
	if err := d.Client.Unscoped().Where("id = ?", ID).First(&alias).Error; err != nil {
		t.Fatalf("load alias %s: %v", ID, err)
	}

	return alias.PausedUntil
}

func TestUpdateAliasPausedUntil(t *testing.T) {
	d := newTestDB(t)
	ctx := context.Background()
	alias := createTestAlias(t, d, "promo@mailx.net", "user-1")
	other := createTestAlias(t, d, "theirs@mailx.net", "user-2")
	pausedUntil := time.Date(2026, 10, 15, 8, 0, 0, 0, time.UTC)

	if err := d.UpdateAliasPausedUntil(ctx, alias.ID, "user-1", &pausedUntil); err != nil {
		t.Fatalf("UpdateAliasPausedUntil() error = %v", err)
	}
	if got := storedPausedUntil(t, d, alias.ID); got == nil || !got.Equal(pausedUntil) {
		t.Errorf("paused_until = %v, want %v", got, pausedUntil)
	}

	// The Edit modal's save leaves the pause alone
	alias.Description = "Giveaways"
	alias.Enabled = false
	if err := d.UpdateAlias(ctx, alias); err != nil {
		t.Fatalf("UpdateAlias() error = %v", err)
	}
	if got := storedPausedUntil(t, d, alias.ID); got == nil || !got.Equal(pausedUntil) {
		t.Errorf("paused_until after UpdateAlias = %v, want %v", got, pausedUntil)
	}

	// Scoped to the owner
	if err := d.UpdateAliasPausedUntil(ctx, other.ID, "user-1", &pausedUntil); err == nil {
		t.Error("UpdateAliasPausedUntil() for another user's alias should fail")
	}
	if got := storedPausedUntil(t, d, other.ID); got != nil {
		t.Errorf("another user's alias was paused until %v", got)
	}
	if err := d.UpdateAliasPausedUntil(ctx, "missing", "user-1", nil); err == nil {
		t.Error("UpdateAliasPausedUntil() for an unknown alias should fail")
	}

	if err := d.UpdateAliasPausedUntil(ctx, alias.ID, "user-1", nil); err != nil {
		t.Fatalf("UpdateAliasPausedUntil(nil) error = %v", err)
	}
	if got := storedPausedUntil(t, d, alias.ID); got != nil {
		t.Errorf("paused_until = %v, want nil", got)
	}

	// Unpausing an alias that isn't paused succeeds
	if err := d.UpdateAliasPausedUntil(ctx, alias.ID, "user-1", nil); err != nil {
		t.Errorf("UpdateAliasPausedUntil(nil) on an unpaused alias error = %v", err)
	}

	// Deleted aliases are excluded
	if err := d.DeleteAlias(ctx, alias.ID, "user-1"); err != nil {
		t.Fatalf("DeleteAlias() error = %v", err)
	}
	if err := d.UpdateAliasPausedUntil(ctx, alias.ID, "user-1", &pausedUntil); err == nil {
		t.Error("UpdateAliasPausedUntil() for a deleted alias should fail")
	}
}
