package service

import (
	"context"
	"errors"

	"github.com/go-webauthn/webauthn/webauthn"
)

var (
	ErrSaveCredential   = errors.New("could not save credential")
	ErrUpdateCredential = errors.New("could not update credential")
	ErrDeleteCredential = errors.New("could not delete credential")
)

type CredentialStore interface {
	SaveCredential(context.Context, webauthn.Credential, string) error
	UpdateCredential(context.Context, webauthn.Credential, string) error
	DeleteCredential(context.Context, webauthn.Credential, string) error
	DeleteCredentialByUserID(context.Context, string) error
	DeleteCredentialByID(context.Context, string, string) error
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
