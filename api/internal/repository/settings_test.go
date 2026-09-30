package repository

import (
	"context"
	"testing"

	"ivpn.net/email/api/internal/model"
)

func TestBlockTrackers_DefaultAndUpdate(t *testing.T) {
	d := newTestDB(t)
	ctx := context.Background()

	// New users get the setting on without doing anything
	if err := d.PostSettings(ctx, model.Settings{UserID: "user-1"}); err != nil {
		t.Fatalf("PostSettings() error = %v", err)
	}

	settings, err := d.GetSettings(ctx, "user-1")
	if err != nil {
		t.Fatalf("GetSettings() error = %v", err)
	}
	if !settings.BlockTrackers {
		t.Error("default block_trackers = false, want true")
	}

	if err := d.UpdateBlockTrackers(ctx, "user-1", false); err != nil {
		t.Fatalf("UpdateBlockTrackers() error = %v", err)
	}
	settings, _ = d.GetSettings(ctx, "user-1")
	if settings.BlockTrackers {
		t.Error("block_trackers = true after turning it off")
	}

	// The full settings update doesn't touch the setting
	settings.FromName = "Me"
	if err := d.UpdateSettings(ctx, settings); err != nil {
		t.Fatalf("UpdateSettings() error = %v", err)
	}
	settings, _ = d.GetSettings(ctx, "user-1")
	if settings.BlockTrackers {
		t.Error("block_trackers turned back on by UpdateSettings")
	}

	if err := d.UpdateBlockTrackers(ctx, "user-1", true); err != nil {
		t.Fatalf("UpdateBlockTrackers() error = %v", err)
	}
	settings, _ = d.GetSettings(ctx, "user-1")
	if !settings.BlockTrackers {
		t.Error("block_trackers = false after turning it on")
	}

	if err := d.UpdateBlockTrackers(ctx, "missing", true); err == nil {
		t.Error("UpdateBlockTrackers() for a user without settings should fail")
	}
}

func TestUpdateRemoveHeaderAndLogIssues(t *testing.T) {
	d := newTestDB(t)
	ctx := context.Background()

	if err := d.PostSettings(ctx, model.Settings{UserID: "user-1"}); err != nil {
		t.Fatalf("PostSettings() error = %v", err)
	}

	for _, want := range []bool{true, false} {
		if err := d.UpdateRemoveHeader(ctx, "user-1", want); err != nil {
			t.Fatalf("UpdateRemoveHeader() error = %v", err)
		}
		if err := d.UpdateLogIssues(ctx, "user-1", !want); err != nil {
			t.Fatalf("UpdateLogIssues() error = %v", err)
		}

		settings, _ := d.GetSettings(ctx, "user-1")
		if settings.RemoveHeader != want {
			t.Errorf("remove_header = %v, want %v", settings.RemoveHeader, want)
		}
		if settings.LogIssues != !want {
			t.Errorf("log_issues = %v, want %v", settings.LogIssues, !want)
		}
		// Each toggle only touches its own column
		if !settings.BlockTrackers {
			t.Error("block_trackers changed by another toggle")
		}
	}

	if err := d.UpdateRemoveHeader(ctx, "missing", true); err == nil {
		t.Error("UpdateRemoveHeader() for a user without settings should fail")
	}
	if err := d.UpdateLogIssues(ctx, "missing", true); err == nil {
		t.Error("UpdateLogIssues() for a user without settings should fail")
	}
}

func TestBlockTrackers_MigrationDefaultsExistingRowsToOn(t *testing.T) {
	d := newTestDB(t)

	// A settings row created before the column existed
	if err := d.Client.Exec(`CREATE TABLE legacy_settings (id text, user_id text)`).Error; err != nil {
		t.Fatalf("create table: %v", err)
	}
	if err := d.Client.Exec(`INSERT INTO legacy_settings (id, user_id) VALUES ('s-1', 'user-1')`).Error; err != nil {
		t.Fatalf("insert: %v", err)
	}
	if err := d.Client.Table("legacy_settings").AutoMigrate(&model.Settings{}); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}

	var settings model.Settings
	if err := d.Client.Table("legacy_settings").Where("user_id = ?", "user-1").First(&settings).Error; err != nil {
		t.Fatalf("First() error = %v", err)
	}
	if !settings.BlockTrackers {
		t.Error("existing row block_trackers = false after migration, want true")
	}
}
