package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetAccessKeys       = errors.New("Unable to retrieve access keys.")
	ErrGetAccessKey        = errors.New("Unable to retrieve access key.")
	ErrPostUserAccessKey   = errors.New("Unable to create access key. Please try again.")
	ErrDeleteUserAccessKey = errors.New("Unable to delete access key. Please try again.")
	ErrInvalidAccessKey    = errors.New("Invalid access key provided.")
	ErrAccessKeyExpired    = errors.New("Access key has expired.")
)

type AccessKeyStore interface {
	GetAccessKeys(context.Context, string) ([]model.AccessKey, error)
	GetAccessKey(context.Context, string) (model.AccessKey, error)
	PostAccessKey(context.Context, model.AccessKey) (model.AccessKey, error)
	DeleteAccessKey(context.Context, string, string) error
	DeleteAccessKeysByUserID(context.Context, string) error
}

func (s *Service) GetAccessKeys(ctx context.Context, userId string) ([]model.AccessKey, error) {
	accessKeys, err := s.Store.GetAccessKeys(ctx, userId)
	if err != nil {
		log.Printf("error getting access keys: %s", err.Error())
		return nil, ErrGetAccessKeys
	}

	return accessKeys, nil
}

func (s *Service) GetAccessKey(ctx context.Context, key string) (model.AccessKey, error) {
	parts := strings.Split(key, ".")
	id := parts[0]
	token := parts[1]

	accessKey, err := s.Store.GetAccessKey(ctx, id)
	if err != nil {
		return model.AccessKey{}, ErrGetAccessKey
	}

	matches := accessKey.Matches(token)
	if !matches {
		return model.AccessKey{}, ErrInvalidAccessKey
	}

	if accessKey.ExpiresAt != nil && accessKey.ExpiresAt.Before(time.Now()) {
		return model.AccessKey{}, ErrAccessKeyExpired
	}

	return accessKey, nil
}

func (s *Service) PostAccessKey(ctx context.Context, userId string, accessKey model.AccessKey) (model.AccessKey, error) {
	if accessKey.TokenPlain != nil {
		err := accessKey.SetToken(*accessKey.TokenPlain)
		if err != nil {
			log.Println("error setting access key token:", err.Error())
			return model.AccessKey{}, ErrPostUserAccessKey
		}
	}

	return s.Store.PostAccessKey(ctx, accessKey)
}

func (s *Service) DeleteAccessKey(ctx context.Context, accessKeyId string, userId string) error {
	err := s.Store.DeleteAccessKey(ctx, accessKeyId, userId)
	if err != nil {
		log.Printf("error deleting access key: %s", err.Error())
		return ErrDeleteUserAccessKey
	}

	return nil
}

func (s *Service) DeleteAccessKeysByUserID(ctx context.Context, userId string) error {
	return s.Store.DeleteAccessKeysByUserID(ctx, userId)
}

func (s *Service) GetDefaults(ctx context.Context, userId string) (model.Settings, string, error) {
	settings, err := s.GetSettings(ctx, userId)
	if err != nil {
		return model.Settings{}, "", err
	}

	rcps, err := s.GetRecipients(ctx, userId)
	if err != nil {
		return model.Settings{}, "", err
	}

	emails := make([]string, 0, len(rcps))
	for _, r := range rcps {
		emails = append(emails, r.Email)
	}
	rcpsStr := strings.Join(emails, ",")

	return settings, rcpsStr, nil
}
