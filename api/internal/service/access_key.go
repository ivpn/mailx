package service

import (
	"context"
	"errors"
	"log"

	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetAccessKeys       = errors.New("Unable to retrieve access keys.")
	ErrPostUserAccessKey   = errors.New("Unable to create access key. Please try again.")
	ErrDeleteUserAccessKey = errors.New("Unable to delete access key. Please try again.")
)

type AccessKeyStore interface {
	GetAccessKeys(context.Context, string) ([]model.AccessKey, error)
	PostAccessKey(context.Context, model.AccessKey) error
	DeleteAccessKey(context.Context, string, string) error
	DeleteAccessKeysByUserID(context.Context, string) error
}

func (s *Service) GetAccessKeys(ctx context.Context, userID string) ([]model.AccessKey, error) {
	accessKeys, err := s.Store.GetAccessKeys(ctx, userID)
	if err != nil {
		log.Printf("error getting access keys: %s", err.Error())
		return nil, ErrGetAccessKeys
	}

	return accessKeys, nil
}

func (s *Service) PostAccessKey(ctx context.Context, accessKey model.AccessKey) error {
	if accessKey.TokenPlain == nil {
		err := accessKey.SetToken(*accessKey.TokenPlain)
		if err != nil {
			log.Println("error setting access key token:", err.Error())
			return ErrPostUserAccessKey
		}
	}

	return s.Store.PostAccessKey(ctx, accessKey)
}

func (s *Service) DeleteAccessKey(ctx context.Context, accessKeyID string, userID string) error {
	err := s.Store.DeleteAccessKey(ctx, accessKeyID, userID)
	if err != nil {
		log.Printf("error deleting access key: %s", err.Error())
		return ErrDeleteUserAccessKey
	}

	return nil
}

func (s *Service) DeleteAccessKeysByUserID(ctx context.Context, userID string) error {
	return s.Store.DeleteAccessKeysByUserID(ctx, userID)
}
