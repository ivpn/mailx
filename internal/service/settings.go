package service

import (
	"context"
	"errors"

	"ivpn.net/email-service/internal/model"
)

var (
	ErrGetSettingsByName = errors.New("could not get settings by user ID")
	ErrPostSettings      = errors.New("could not post settings")
	ErrUpdateSettings    = errors.New("could not update settings")
	ErrDeleteSettings    = errors.New("could not delete settings")
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
		return model.Settings{}, err
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
