package service

import (
	"context"
	"errors"

	"github.com/go-webauthn/webauthn/webauthn"
	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetSession    = errors.New("could not get session by token")
	ErrSaveSession   = errors.New("could not save session")
	ErrDeleteSession = errors.New("could not delete session")
)

type SessionStore interface {
	GetSession(context.Context, string) (model.Session, bool, error)
	SaveSession(context.Context, webauthn.SessionData, string, string) error
	DeleteSession(context.Context, string) error
}

func (s *Service) GetSession(ctx context.Context, token string) (model.Session, bool, error) {
	session, exists, err := s.Store.GetSession(ctx, token)
	if err != nil {
		return model.Session{}, false, ErrGetSession
	}

	return session, exists, nil
}

func (s *Service) SaveSession(ctx context.Context, session webauthn.SessionData, token string, userID string) error {
	err := s.Store.SaveSession(ctx, session, token, userID)
	if err != nil {
		return ErrSaveSession
	}

	return nil
}

func (s *Service) DeleteSession(ctx context.Context, token string) error {
	err := s.Store.DeleteSession(ctx, token)
	if err != nil {
		return ErrDeleteSession
	}

	return nil
}
