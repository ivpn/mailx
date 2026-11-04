package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"ivpn.net/email/api/internal/model"
)

func newSessionTestDB(t *testing.T) *Database {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&model.Session{}, &model.AccessKey{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return &Database{Client: db}
}

func TestGetSessionRejectsExpired(t *testing.T) {
	ctx := context.Background()
	d := newSessionTestDB(t)

	err := d.SaveSession(ctx, model.Session{UserID: "user-1", Token: "valid", ExpiresAt: time.Now().Add(time.Hour), Scope: model.SessionScopeWeb})
	assert.NoError(t, err)
	err = d.SaveSession(ctx, model.Session{UserID: "user-1", Token: "expired", ExpiresAt: time.Now().Add(-time.Second), Scope: model.SessionScopeWeb})
	assert.NoError(t, err)

	session, ok, err := d.GetSession(ctx, "valid")
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "user-1", session.UserID)

	_, ok, err = d.GetSession(ctx, "expired")
	assert.Error(t, err)
	assert.False(t, ok)
}

func TestSaveSessionStoresScope(t *testing.T) {
	ctx := context.Background()
	d := newSessionTestDB(t)
	exp := time.Now().Add(time.Hour)

	key, err := d.PostAccessKey(ctx, model.AccessKey{UserId: "user-1", Name: "extension"})
	require.NoError(t, err)

	require.NoError(t, d.SaveSession(ctx, model.Session{UserID: "user-1", Token: "web", ExpiresAt: exp, Scope: model.SessionScopeWeb}))
	require.NoError(t, d.SaveSession(ctx, model.Session{UserID: "user-1", Token: "api", ExpiresAt: exp, Scope: model.SessionScopeAPI, AccessKeyID: key.ID}))

	web, _, err := d.GetSession(ctx, "web")
	require.NoError(t, err)
	assert.Equal(t, model.SessionScopeWeb, web.Scope)
	assert.Empty(t, web.AccessKeyID)

	api, _, err := d.GetSession(ctx, "api")
	require.NoError(t, err)
	assert.Equal(t, model.SessionScopeAPI, api.Scope)
	assert.Equal(t, key.ID, api.AccessKeyID)
}

func TestSaveSessionRequiresExistingAccessKey(t *testing.T) {
	ctx := context.Background()
	d := newSessionTestDB(t)
	exp := time.Now().Add(time.Hour)

	key, err := d.PostAccessKey(ctx, model.AccessKey{UserId: "user-1", Name: "extension"})
	require.NoError(t, err)
	other, err := d.PostAccessKey(ctx, model.AccessKey{UserId: "user-2", Name: "extension"})
	require.NoError(t, err)

	// A key deleted after it was checked at login, but before the session was saved
	require.NoError(t, d.DeleteAccessKey(ctx, key.ID, "user-1"))
	err = d.SaveSession(ctx, model.Session{UserID: "user-1", Token: "deleted-key", ExpiresAt: exp, Scope: model.SessionScopeAPI, AccessKeyID: key.ID})
	assert.ErrorIs(t, err, ErrAccessKeyNotFound)

	err = d.SaveSession(ctx, model.Session{UserID: "user-1", Token: "other-user-key", ExpiresAt: exp, Scope: model.SessionScopeAPI, AccessKeyID: other.ID})
	assert.ErrorIs(t, err, ErrAccessKeyNotFound)

	var count int64
	require.NoError(t, d.Client.Model(&model.Session{}).Count(&count).Error)
	assert.Zero(t, count, "no session is saved without its access key")
}

func TestSaveSessionRequiresScope(t *testing.T) {
	ctx := context.Background()
	d := newSessionTestDB(t)

	for _, scope := range []model.SessionScope{"", "admin"} {
		err := d.SaveSession(ctx, model.Session{UserID: "user-1", Token: "token-" + string(scope), ExpiresAt: time.Now().Add(time.Hour), Scope: scope})
		assert.ErrorIs(t, err, model.ErrInvalidSessionScope, "scope %q", scope)
	}

	var count int64
	require.NoError(t, d.Client.Model(&model.Session{}).Count(&count).Error)
	assert.Zero(t, count, "no session is saved without a valid scope")
}

func TestBackfillSessionScopes(t *testing.T) {
	ctx := context.Background()
	d := newSessionTestDB(t)

	const tokenExpiration = 168 * time.Hour
	created := time.Now().Add(-time.Hour)

	// Sessions as the old code saved them: no scope, expiry computed just before the row is saved
	legacy := []struct {
		token     string
		lifetime  time.Duration
		remember  bool
		wantScope model.SessionScope
	}{
		{token: "web", lifetime: tokenExpiration - time.Millisecond, wantScope: model.SessionScopeWeb},
		{token: "web-remember", lifetime: 2160 * time.Hour, remember: true, wantScope: model.SessionScopeWeb},
		{token: "api", lifetime: 8760 * time.Hour, wantScope: model.SessionScopeAPI},
		{token: "api-key-expiring-in-a-month", lifetime: 30 * 24 * time.Hour, wantScope: model.SessionScopeAPI},
		// Can't be told apart from a web session; it expires within tokenExpiration
		{token: "api-key-expiring-in-a-day", lifetime: 24 * time.Hour, wantScope: model.SessionScopeWeb},
		{token: "expired-api", lifetime: 8760 * time.Hour, wantScope: model.SessionScopeAPI},
	}
	for _, s := range legacy {
		require.NoError(t, d.Client.Create(&model.Session{
			BaseModel: model.BaseModel{CreatedAt: created},
			UserID:    "user-1",
			Token:     s.token,
			ExpiresAt: created.Add(s.lifetime),
			Remember:  s.remember,
		}).Error)
	}
	// A column added by AutoMigrate is NULL on existing rows, not ""
	require.NoError(t, d.Client.Model(&model.Session{}).Where("token = ?", "expired-api").UpdateColumn("scope", gorm.Expr("NULL")).Error)

	// Saved by the new code: left alone
	require.NoError(t, d.SaveSession(ctx, model.Session{UserID: "user-1", Token: "new-web", ExpiresAt: time.Now().Add(8760 * time.Hour), Scope: model.SessionScopeWeb}))

	n, err := d.BackfillSessionScopes(ctx, tokenExpiration)
	require.NoError(t, err)
	assert.Equal(t, len(legacy), n)

	for _, s := range legacy {
		var got model.Session
		require.NoError(t, d.Client.Where("token = ?", s.token).First(&got).Error)
		assert.Equal(t, s.wantScope, got.Scope, s.token)
	}

	var newWeb model.Session
	require.NoError(t, d.Client.Where("token = ?", "new-web").First(&newWeb).Error)
	assert.Equal(t, model.SessionScopeWeb, newWeb.Scope)

	n, err = d.BackfillSessionScopes(ctx, tokenExpiration)
	require.NoError(t, err)
	assert.Zero(t, n, "a second run finds nothing to do")
}

func TestDeleteAccessKeyEndsItsSessions(t *testing.T) {
	ctx := context.Background()
	d := newSessionTestDB(t)
	exp := time.Now().Add(time.Hour)

	key, err := d.PostAccessKey(ctx, model.AccessKey{UserId: "user-1", Name: "extension"})
	require.NoError(t, err)
	other, err := d.PostAccessKey(ctx, model.AccessKey{UserId: "user-1", Name: "script"})
	require.NoError(t, err)

	require.NoError(t, d.SaveSession(ctx, model.Session{UserID: "user-1", Token: "from-key", ExpiresAt: exp, Scope: model.SessionScopeAPI, AccessKeyID: key.ID}))
	require.NoError(t, d.SaveSession(ctx, model.Session{UserID: "user-1", Token: "from-other-key", ExpiresAt: exp, Scope: model.SessionScopeAPI, AccessKeyID: other.ID}))
	require.NoError(t, d.SaveSession(ctx, model.Session{UserID: "user-1", Token: "web", ExpiresAt: exp, Scope: model.SessionScopeWeb}))
	// Backfilled: api, but without the key it came from (NULL, as AutoMigrate leaves the new column)
	require.NoError(t, d.SaveSession(ctx, model.Session{UserID: "user-1", Token: "backfilled-api", ExpiresAt: exp, Scope: model.SessionScopeAPI}))
	require.NoError(t, d.Client.Model(&model.Session{}).Where("token = ?", "backfilled-api").UpdateColumn("access_key_id", gorm.Expr("NULL")).Error)

	// Another user can't end the key's sessions, or delete the key
	require.NoError(t, d.DeleteAccessKey(ctx, key.ID, "user-2"))
	_, ok, _ := d.GetSession(ctx, "from-key")
	assert.True(t, ok)
	_, ok, _ = d.GetSession(ctx, "backfilled-api")
	assert.True(t, ok)
	keys, err := d.GetAccessKeys(ctx, "user-1")
	require.NoError(t, err)
	assert.Len(t, keys, 2)

	require.NoError(t, d.SaveSession(ctx, model.Session{UserID: "user-2", Token: "other-user-backfilled-api", ExpiresAt: exp, Scope: model.SessionScopeAPI}))

	require.NoError(t, d.DeleteAccessKey(ctx, key.ID, "user-1"))

	_, ok, _ = d.GetSession(ctx, "from-key")
	assert.False(t, ok, "the key's session ended")
	_, ok, _ = d.GetSession(ctx, "backfilled-api")
	assert.False(t, ok, "a backfilled api session, which may be the key's, ended")
	_, ok, _ = d.GetSession(ctx, "from-other-key")
	assert.True(t, ok, "another key's session is kept")
	_, ok, _ = d.GetSession(ctx, "web")
	assert.True(t, ok, "the web session is kept")
	_, ok, _ = d.GetSession(ctx, "other-user-backfilled-api")
	assert.True(t, ok, "another user's session is kept")
}

func TestDeleteAccessKeysByUserIDEndsApiSessions(t *testing.T) {
	ctx := context.Background()
	d := newSessionTestDB(t)
	exp := time.Now().Add(time.Hour)

	key, err := d.PostAccessKey(ctx, model.AccessKey{UserId: "user-1", Name: "extension"})
	require.NoError(t, err)

	require.NoError(t, d.SaveSession(ctx, model.Session{UserID: "user-1", Token: "from-key", ExpiresAt: exp, Scope: model.SessionScopeAPI, AccessKeyID: key.ID}))
	// Backfilled: api, but without the key it came from
	require.NoError(t, d.SaveSession(ctx, model.Session{UserID: "user-1", Token: "backfilled-api", ExpiresAt: exp, Scope: model.SessionScopeAPI}))
	require.NoError(t, d.SaveSession(ctx, model.Session{UserID: "user-1", Token: "web", ExpiresAt: exp, Scope: model.SessionScopeWeb}))
	require.NoError(t, d.SaveSession(ctx, model.Session{UserID: "user-2", Token: "other-user-api", ExpiresAt: exp, Scope: model.SessionScopeAPI}))

	require.NoError(t, d.DeleteAccessKeysByUserID(ctx, "user-1"))

	_, ok, _ := d.GetSession(ctx, "from-key")
	assert.False(t, ok)
	_, ok, _ = d.GetSession(ctx, "backfilled-api")
	assert.False(t, ok)
	_, ok, _ = d.GetSession(ctx, "web")
	assert.True(t, ok, "the web session is kept")
	_, ok, _ = d.GetSession(ctx, "other-user-api")
	assert.True(t, ok, "another user's session is kept")

	keys, err := d.GetAccessKeys(ctx, "user-1")
	require.NoError(t, err)
	assert.Empty(t, keys)
}
