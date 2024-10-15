package service

import (
	"context"
	"errors"

	"github.com/go-webauthn/webauthn/webauthn"
	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetCredentials   = errors.New("could not get credentials")
	ErrSaveCredential   = errors.New("could not save credential")
	ErrUpdateCredential = errors.New("could not update credential")
	ErrDeleteCredential = errors.New("could not delete credential")
)

type CredentialStore interface {
	GetCredentials(context.Context, string) ([]model.Credential, error)
	SaveCredential(context.Context, webauthn.Credential, string) error
	UpdateCredential(context.Context, webauthn.Credential, string) error
	DeleteCredential(context.Context, webauthn.Credential, string) error
	DeleteCredentialByUserID(context.Context, string) error
	DeleteCredentialByID(context.Context, string, string) error
}

func (s *Service) GetCredentials(ctx context.Context, userID string) ([]model.Credential, error) {
	credentials, err := s.Store.GetCredentials(ctx, userID)
	if err != nil {
		return nil, ErrGetCredentials
	}

	return credentials, nil
}

func (s *Service) SaveCredential(ctx context.Context, credential webauthn.Credential, userID string) error {
	err := s.Store.SaveCredential(ctx, credential, userID)
	if err != nil {
		return ErrSaveCredential
	}

	return nil
}

func (s *Service) UpdateCredential(ctx context.Context, credential webauthn.Credential, userID string) error {
	err := s.Store.UpdateCredential(ctx, credential, userID)
	if err != nil {
		return ErrUpdateCredential
	}

	return nil
}

func (s *Service) DeleteCredential(ctx context.Context, credential webauthn.Credential, userID string) error {
	err := s.Store.DeleteCredential(ctx, credential, userID)
	if err != nil {
		return ErrDeleteCredential
	}

	return nil
}

func (s *Service) DeleteCredentialByID(ctx context.Context, ID string, userID string) error {
	err := s.Store.DeleteCredentialByID(ctx, ID, userID)
	if err != nil {
		return ErrDeleteCredential
	}

	return nil
}
