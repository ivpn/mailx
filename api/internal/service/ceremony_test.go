package service

import (
	"context"
	"testing"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/stretchr/testify/assert"
	"ivpn.net/email/api/internal/model"
)

var allCeremonyTypes = []model.CeremonyType{
	model.CeremonyRegistration,
	model.CeremonyLogin,
	model.CeremonyPasskeyLogin,
	model.CeremonyAddPasskey,
	model.CeremonyStepUp,
}

func TestCeremonyIsSingleUse(t *testing.T) {
	ctx := context.Background()
	cache := newTotpTestCache()
	s := &Service{Cache: cache}

	ceremony := model.Ceremony{
		Type:        model.CeremonyLogin,
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

	got, err := s.ConsumeCeremony(ctx, token, model.CeremonyLogin)
	assert.NoError(t, err)
	assert.Equal(t, model.CeremonyLogin, got.Type)
	assert.Equal(t, "user-1", got.UserID)
	assert.Equal(t, "challenge", got.SessionData.Challenge)
	assert.True(t, got.Remember)
	assert.WithinDuration(t, time.Now().Add(5*time.Minute), got.SessionData.Expires, 5*time.Second)

	_, err = s.ConsumeCeremony(ctx, token, model.CeremonyLogin)
	assert.ErrorIs(t, err, ErrGetCeremony)
}

func TestConsumeCeremonyRejectsUnknownToken(t *testing.T) {
	s := &Service{Cache: newTotpTestCache()}

	_, err := s.ConsumeCeremony(context.Background(), "", model.CeremonyLogin)
	assert.ErrorIs(t, err, ErrGetCeremony)

	_, err = s.ConsumeCeremony(context.Background(), "unknown", model.CeremonyLogin)
	assert.ErrorIs(t, err, ErrGetCeremony)
}

func TestConsumeCeremonyRejectsOtherType(t *testing.T) {
	ctx := context.Background()

	for _, saved := range allCeremonyTypes {
		for _, want := range allCeremonyTypes {
			if saved == want {
				continue
			}

			t.Run(string(saved)+"->"+string(want), func(t *testing.T) {
				s := &Service{Cache: newTotpTestCache()}
				token, err := s.SaveCeremony(ctx, model.Ceremony{Type: saved, UserID: "user-1"}, 5*time.Minute)
				assert.NoError(t, err)

				_, err = s.ConsumeCeremony(ctx, token, want)
				assert.ErrorIs(t, err, ErrGetCeremony)

				// The rejected ceremony is gone, so it can't be retried at its own finish
				_, err = s.ConsumeCeremony(ctx, token, saved)
				assert.ErrorIs(t, err, ErrGetCeremony)
			})
		}
	}
}

func TestSaveCeremonyRequiresType(t *testing.T) {
	s := &Service{Cache: newTotpTestCache()}

	_, err := s.SaveCeremony(context.Background(), model.Ceremony{UserID: "user-1"}, 5*time.Minute)
	assert.ErrorIs(t, err, ErrSaveCeremony)
}

func TestConsumeCeremonyRejectsUntypedState(t *testing.T) {
	ctx := context.Background()
	cache := newTotpTestCache()
	s := &Service{Cache: cache}

	// State saved before ceremonies were typed has no "type" field
	err := cache.Set(ctx, ceremonyCacheKey("legacy"), `{"user_id":"user-1","session_data":{"challenge":"c"}}`, time.Minute)
	assert.NoError(t, err)

	for _, want := range allCeremonyTypes {
		_, err = s.ConsumeCeremony(ctx, "legacy", want)
		assert.ErrorIs(t, err, ErrGetCeremony)
	}
}
