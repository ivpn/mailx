package service

import (
	"context"
	"errors"
	"log"

	"github.com/go-webauthn/webauthn/webauthn"
	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetCredentials        = errors.New("Unable to retrieve credentials.")
	ErrSaveCredential        = errors.New("Unable to save credential. Please try again.")
	ErrUpdateCredential      = errors.New("Unable to update credential. Please try again.")
	ErrDeleteCredential      = errors.New("Unable to delete credential. Please try again.")
	ErrMaxExceededCredential = errors.New("You have reached the maximum number of allowed passkeys.")
)

type CredentialStore interface {
	GetCredentials(context.Context, string) ([]model.Credential, error)
	GetCredentialsCount(context.Context, string) (int, error)
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
	count, err := s.Store.GetCredentialsCount(ctx, userID)
	if err != nil {
		log.Printf("error saving credential: %s", err.Error())
		return ErrSaveCredential
	}

	if count >= s.Cfg.Service.MaxCredentials {
		return ErrMaxExceededCredential
	}

	err = s.Store.SaveCredential(ctx, credential, userID)
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
