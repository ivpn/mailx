package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"ivpn.net/email/api/internal/model"
)

var (
	ErrSaveCeremony = errors.New("Unable to save session. Please try again.")
	ErrGetCeremony  = errors.New("Unable to retrieve session. Please try again.")
)

// Ceremony state lives only in the cache under its own namespace, never in the sessions table,
// so a ceremony token can't be used as a login session.
func ceremonyCacheKey(token string) string {
	return "webauthn_ceremony_" + token
}

// SaveCeremony stores the ceremony under a new random token and returns the token. The WebAuthn
// session data expiry is set to match the cache TTL, so the library enforces it on finish too.
func (s *Service) SaveCeremony(ctx context.Context, ceremony model.Ceremony, ttl time.Duration) (string, error) {
	token, err := model.GenSessionToken()
	if err != nil {
		log.Printf("error generating ceremony token: %s", err.Error())
		return "", ErrSaveCeremony
	}

	ceremony.SessionData.Expires = time.Now().Add(ttl)
	data, err := json.Marshal(ceremony)
	if err != nil {
		log.Printf("error marshalling ceremony: %s", err.Error())
		return "", ErrSaveCeremony
	}

	err = s.Cache.Set(ctx, ceremonyCacheKey(token), string(data), ttl)
	if err != nil {
		log.Printf("error saving ceremony: %s", err.Error())
		return "", ErrSaveCeremony
	}

	return token, nil
}

// ConsumeCeremony is single-use: the ceremony is removed as it is read, so a challenge can't be
// replayed, whether or not the finish step succeeds.
func (s *Service) ConsumeCeremony(ctx context.Context, token string) (model.Ceremony, error) {
	if token == "" {
		return model.Ceremony{}, ErrGetCeremony
	}

	data, err := s.Cache.GetDel(ctx, ceremonyCacheKey(token))
	if err != nil {
		return model.Ceremony{}, ErrGetCeremony
	}

	var ceremony model.Ceremony
	err = json.Unmarshal([]byte(data), &ceremony)
	if err != nil {
		log.Printf("error unmarshalling ceremony: %s", err.Error())
		return model.Ceremony{}, ErrGetCeremony
	}

	return ceremony, nil
}
