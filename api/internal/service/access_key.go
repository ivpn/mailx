package service

import (
	"context"
	"errors"
	"log"
	"strings"

	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetAccessKeys       = errors.New("Unable to retrieve access keys.")
	ErrGetAccessKey        = errors.New("Unable to retrieve access key.")
	ErrPostUserAccessKey   = errors.New("Unable to create access key. Please try again.")
	ErrDeleteUserAccessKey = errors.New("Unable to delete access key. Please try again.")
	ErrInvalidAccessKey    = errors.New("Invalid access key provided.")
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
		log.Printf("error getting access key by hash: %s", err.Error())
		return model.AccessKey{}, ErrGetAccessKey
	}

	matches := accessKey.Matches(token)
	if !matches {
		log.Printf("access key token does not match: %s", key)
		return model.AccessKey{}, ErrInvalidAccessKey
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
