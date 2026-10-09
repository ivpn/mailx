package service

import (
	"context"
	"testing"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/stretchr/testify/assert"
	"ivpn.net/email/api/internal/model"
)

func TestCeremonyIsSingleUse(t *testing.T) {
	ctx := context.Background()
	cache := newTotpTestCache()
	s := &Service{Cache: cache}

	ceremony := model.Ceremony{
		UserID:      "user-1",
		SessionData: webauthn.SessionData{Challenge: "challenge"},
		Remember:    true,
	}
	token, err := s.SaveCeremony(ctx, ceremony, 5*time.Minute)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Stored only under the ceremony namespace, never as a bare token
	_, err = cache.Get(ctx, token)
	assert.Error(t, err)

	got, err := s.ConsumeCeremony(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, "user-1", got.UserID)
	assert.Equal(t, "challenge", got.SessionData.Challenge)
	assert.True(t, got.Remember)
	assert.WithinDuration(t, time.Now().Add(5*time.Minute), got.SessionData.Expires, 5*time.Second)

	_, err = s.ConsumeCeremony(ctx, token)
	assert.ErrorIs(t, err, ErrGetCeremony)
}

func TestConsumeCeremonyRejectsUnknownToken(t *testing.T) {
	s := &Service{Cache: newTotpTestCache()}

	_, err := s.ConsumeCeremony(context.Background(), "")
	assert.ErrorIs(t, err, ErrGetCeremony)

	_, err = s.ConsumeCeremony(context.Background(), "unknown")
	assert.ErrorIs(t, err, ErrGetCeremony)
}
