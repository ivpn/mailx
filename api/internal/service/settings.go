package service

import (
	"context"
	"errors"

	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetSettings    = errors.New("Unable to retrieve settings by user ID.")
	ErrPostSettings   = errors.New("Unable to create settings.")
	ErrUpdateSettings = errors.New("Unable to update settings.")
	ErrDeleteSettings = errors.New("Unable to delete settings.")
)

type SettingsStore interface {
	GetSettings(context.Context, string) (model.Settings, error)
	PostSettings(context.Context, model.Settings) error
	UpdateSettings(context.Context, model.Settings) error
	DeleteSettings(context.Context, string) error
}

func (s *Service) GetSettings(ctx context.Context, userID string) (model.Settings, error) {
	settings, err := s.Store.GetSettings(ctx, userID)
	if err != nil {
		return model.Settings{}, ErrGetSettings
	}

	return settings, nil
}

func (s *Service) PostSettings(ctx context.Context, userID string) error {
	settings := model.Settings{
		UserID: userID,
	}

	err := s.Store.PostSettings(ctx, settings)
	if err != nil {
		return ErrPostSettings
	}

	return nil
}

func (s *Service) UpdateSettings(ctx context.Context, settings model.Settings) error {
	err := s.Store.UpdateSettings(ctx, settings)
	if err != nil {
		return ErrUpdateSettings
	}

	return nil
}

func (s *Service) DeleteSettings(ctx context.Context, userID string) error {
	err := s.Store.DeleteSettings(ctx, userID)
	if err != nil {
		return ErrDeleteSettings
	}

	return nil
}
