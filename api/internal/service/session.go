package service

import (
	"context"
	"errors"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetSession    = errors.New("Unable to retrieve session by token.")
	ErrSaveSession   = errors.New("Unable to save session. Please try again.")
	ErrDeleteSession = errors.New("Unable to delete session. Please try again.")
)

type SessionStore interface {
	GetSession(context.Context, string) (model.Session, bool, error)
	GetSessionCount(context.Context, string) (int, error)
	SaveSession(context.Context, webauthn.SessionData, string, string, time.Time) error
	DeleteSession(context.Context, string) error
	DeleteSessionByUserID(context.Context, string) error
}

func (s *Service) GetSession(ctx context.Context, token string) (model.Session, bool, error) {
	session, exists, err := s.Store.GetSession(ctx, token)
	if err != nil {
		return model.Session{}, false, ErrGetSession
	}

	return session, exists, nil
}

func (s *Service) GetSessionCount(ctx context.Context, userID string) (int, error) {
	count, err := s.Store.GetSessionCount(ctx, userID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Service) CheckSessionCount(ctx context.Context, userID string) (bool, error) {
	count, err := s.Store.GetSessionCount(ctx, userID)
	if err != nil {
		return false, err
	}

	if count >= s.Cfg.Service.MaxSessions {
		return false, nil
	}

	return true, nil
}

func (s *Service) SaveSession(ctx context.Context, session webauthn.SessionData, token string, userID string, exp time.Time) error {
	err := s.Store.SaveSession(ctx, session, token, userID, exp)
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
