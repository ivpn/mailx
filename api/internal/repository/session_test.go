package repository

import (
	"context"
	"testing"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"ivpn.net/email/api/internal/model"
)

func newSessionTestDB(t *testing.T) *Database {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&model.Session{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return &Database{Client: db}
}

func TestGetSessionRejectsExpired(t *testing.T) {
	ctx := context.Background()
	d := newSessionTestDB(t)

	err := d.SaveSession(ctx, webauthn.SessionData{}, "valid", "user-1", time.Now().Add(time.Hour), false)
	assert.NoError(t, err)
	err = d.SaveSession(ctx, webauthn.SessionData{}, "expired", "user-1", time.Now().Add(-time.Second), false)
	assert.NoError(t, err)

	session, ok, err := d.GetSession(ctx, "valid")
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "user-1", session.UserID)

	_, ok, err = d.GetSession(ctx, "expired")
	assert.Error(t, err)
	assert.False(t, ok)
}
