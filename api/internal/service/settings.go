package service

import (
	"context"
	"errors"

	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetSettings         = errors.New("Unable to retrieve settings by user ID.")
	ErrPostSettings        = errors.New("Unable to create settings.")
	ErrUpdateSettings      = errors.New("Unable to update settings.")
	ErrDeleteSettings      = errors.New("Unable to delete settings.")
	ErrUpdateBlockTrackers = errors.New("Unable to update block email trackers setting.")
	ErrUpdateRemoveHeader  = errors.New("Unable to update Mailx header setting.")
	ErrUpdateLogIssues     = errors.New("Unable to update diagnostics setting.")
	ErrUpdateDefaultDomain = errors.New("Unable to update default domain.")
	ErrUpdateDefaultRcp    = errors.New("Unable to update default recipient.")
	ErrUpdateAliasFormat   = errors.New("Unable to update default alias format.")
	ErrUpdateFromName      = errors.New("Unable to update from name.")
)

type SettingsStore interface {
	GetSettings(context.Context, string) (model.Settings, error)
	PostSettings(context.Context, model.Settings) error
	UpdateSettings(context.Context, model.Settings) error
	UpdateListUnsubscribeBehaviour(context.Context, string, model.ListUnsubscribeBehaviour) error
	UpdateBlockTrackers(context.Context, string, bool) error
	UpdateRemoveHeader(context.Context, string, bool) error
	UpdateLogIssues(context.Context, string, bool) error
	UpdateDefaultDomain(context.Context, string, string) error
	UpdateDefaultRecipient(context.Context, string, string) error
	UpdateAliasFormat(context.Context, string, string) error
	UpdateFromName(context.Context, string, string) error
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

func (s *Service) UpdateBlockTrackers(ctx context.Context, userID string, enabled bool) error {
	err := s.Store.UpdateBlockTrackers(ctx, userID, enabled)
	if err != nil {
		return ErrUpdateBlockTrackers
	}

	return nil
}

func (s *Service) UpdateRemoveHeader(ctx context.Context, userID string, removeHeader bool) error {
	err := s.Store.UpdateRemoveHeader(ctx, userID, removeHeader)
	if err != nil {
		return ErrUpdateRemoveHeader
	}

	return nil
}

func (s *Service) UpdateLogIssues(ctx context.Context, userID string, enabled bool) error {
	err := s.Store.UpdateLogIssues(ctx, userID, enabled)
	if err != nil {
		return ErrUpdateLogIssues
	}

	return nil
}

func (s *Service) UpdateDefaultDomain(ctx context.Context, userID string, domain string) error {
	err := s.Store.UpdateDefaultDomain(ctx, userID, domain)
	if err != nil {
		return ErrUpdateDefaultDomain
	}

	return nil
}

func (s *Service) UpdateDefaultRecipient(ctx context.Context, userID string, recipient string) error {
	err := s.Store.UpdateDefaultRecipient(ctx, userID, recipient)
	if err != nil {
		return ErrUpdateDefaultRcp
	}

	return nil
}

func (s *Service) UpdateAliasFormat(ctx context.Context, userID string, format string) error {
	err := s.Store.UpdateAliasFormat(ctx, userID, format)
	if err != nil {
		return ErrUpdateAliasFormat
	}

	return nil
}

func (s *Service) UpdateFromName(ctx context.Context, userID string, fromName string) error {
	err := s.Store.UpdateFromName(ctx, userID, fromName)
	if err != nil {
		return ErrUpdateFromName
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
