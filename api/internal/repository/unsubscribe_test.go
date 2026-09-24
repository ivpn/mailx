package repository

import (
	"context"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"ivpn.net/email/api/internal/model"
)

func newTestDB(t *testing.T) *Database {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	if err := db.AutoMigrate(&model.Settings{}, &model.Alias{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return &Database{Client: db}
}

func TestListUnsubscribeBehaviour_DefaultAndUpdate(t *testing.T) {
	d := newTestDB(t)
	ctx := context.Background()

	if err := d.PostSettings(ctx, model.Settings{UserID: "user-1"}); err != nil {
		t.Fatalf("PostSettings() error = %v", err)
	}

	settings, err := d.GetSettings(ctx, "user-1")
	if err != nil {
		t.Fatalf("GetSettings() error = %v", err)
	}
	if settings.ListUnsubscribeBehaviour != model.ListUnsubscribeOriginalWithFallback {
		t.Errorf("default behaviour = %d, want OriginalWithFallback", settings.ListUnsubscribeBehaviour)
	}

	if err := d.UpdateListUnsubscribeBehaviour(ctx, "user-1", model.ListUnsubscribeBlockEmail); err != nil {
		t.Fatalf("UpdateListUnsubscribeBehaviour() error = %v", err)
	}
	settings, _ = d.GetSettings(ctx, "user-1")
	if settings.ListUnsubscribeBehaviour != model.ListUnsubscribeBlockEmail {
		t.Errorf("behaviour = %d, want BlockEmail", settings.ListUnsubscribeBehaviour)
	}

	// The full settings update doesn't touch the behaviour
	settings.FromName = "Me"
	if err := d.UpdateSettings(ctx, settings); err != nil {
		t.Fatalf("UpdateSettings() error = %v", err)
	}
	settings, _ = d.GetSettings(ctx, "user-1")
	if settings.ListUnsubscribeBehaviour != model.ListUnsubscribeBlockEmail {
		t.Errorf("behaviour after UpdateSettings = %d, want BlockEmail", settings.ListUnsubscribeBehaviour)
	}

	if err := d.UpdateListUnsubscribeBehaviour(ctx, "missing", model.ListUnsubscribeDelete); err == nil {
		t.Error("UpdateListUnsubscribeBehaviour() for a user without settings should fail")
	}
}

func TestOneClickAliasLookups(t *testing.T) {
	d := newTestDB(t)
	ctx := context.Background()

	alias := model.Alias{Name: "gone@domain.net", UserID: "user-1", Enabled: true}
	if err := d.Client.Create(&alias).Error; err != nil {
		t.Fatalf("create alias error = %v", err)
	}

	if deleted, err := d.IsAliasNameDeleted("gone@domain.net"); err != nil || deleted {
		t.Errorf("IsAliasNameDeleted() before delete = %v, %v, want false", deleted, err)
	}

	if err := d.DeleteAlias(ctx, alias.ID, "user-1"); err != nil {
		t.Fatalf("DeleteAlias() error = %v", err)
	}

	if deleted, err := d.IsAliasNameDeleted("gone@domain.net"); err != nil || !deleted {
		t.Errorf("IsAliasNameDeleted() after delete = %v, %v, want true", deleted, err)
	}
	if deleted, _ := d.IsAliasNameDeleted("other@domain.net"); deleted {
		t.Error("IsAliasNameDeleted() for an unknown name = true")
	}

	got, err := d.GetAliasByIDUnscoped(ctx, alias.ID)
	if err != nil || got.ID != alias.ID || !got.DeletedAt.Valid {
		t.Errorf("GetAliasByIDUnscoped() = %+v, %v, want the deleted alias", got, err)
	}
	if _, err := d.GetAliasByIDUnscoped(ctx, "missing"); err == nil {
		t.Error("GetAliasByIDUnscoped() for an unknown ID should fail")
	}
}
