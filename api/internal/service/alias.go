package service

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetAlias            = errors.New("could not get alias by ID")
	ErrGetAliases          = errors.New("could not get aliass by recipient ID")
	ErrGetAliasByName      = errors.New("could not get alias by name")
	ErrPostAlias           = errors.New("could not create alias, try again")
	ErrUpdateAlias         = errors.New("could not update alias")
	ErrDeleteAlias         = errors.New("could not delete alias")
	ErrDeleteAliasByUserID = errors.New("could not delete alias by user ID")
)

type AliasStore interface {
	GetAlias(context.Context, string) (model.Alias, error)
	GetAliases(context.Context, string) ([]model.Alias, error)
	GetAliasByName(string) (model.Alias, error)
	PostAlias(context.Context, model.Alias) error
	UpdateAlias(context.Context, model.Alias) error
	DeleteAlias(context.Context, string) error
	DeleteAliasByUserID(context.Context, string) error
}

func (s *Service) GetAlias(ctx context.Context, ID string) (model.Alias, error) {
	alias, err := s.Store.GetAlias(ctx, ID)
	if err != nil {
		log.Printf("error fetching alias: %s", err.Error())
		return model.Alias{}, ErrGetAlias
	}

	return alias, nil
}

func (s *Service) GetAliases(ctx context.Context, userID string) ([]model.Alias, error) {
	aliases, err := s.Store.GetAliases(ctx, userID)
	if err != nil {
		log.Printf("error fetching aliass: %s", err.Error())
		return []model.Alias{}, ErrGetAliases
	}

	return aliases, nil
}

func (s *Service) GetAliasByName(name string) (model.Alias, error) {
	alias, err := s.Store.GetAliasByName(name)
	if err != nil {
		log.Printf("error fetching alias: %s", err.Error())
		return model.Alias{}, ErrGetAliasByName
	}

	return alias, nil
}

func (s *Service) PostAlias(ctx context.Context, alias model.Alias, format string, domain string) error {
	for i := 0; i < 5; i++ {
		alias.Name = model.GenerateAlias(format) + "@" + domain
		err := s.Store.PostAlias(ctx, alias)
		if err != nil {
			log.Printf("error creating alias: %s", err.Error())
			switch {
			case errors.Is(err, gorm.ErrDuplicatedKey):
				continue
			default:
				return ErrPostAlias
			}
		}
		break
	}

	return nil
}

func (s *Service) UpdateAlias(ctx context.Context, alias model.Alias) error {
	err := s.Store.UpdateAlias(ctx, alias)
	if err != nil {
		log.Printf("error updating alias: %s", err.Error())
		return ErrUpdateAlias
	}

	return nil
}

func (s *Service) DeleteAlias(ctx context.Context, ID string) error {
	err := s.Store.DeleteAlias(ctx, ID)
	if err != nil {
		log.Printf("error deleting alias: %s", err.Error())
		return ErrDeleteAlias
	}

	return nil
}

func (s *Service) DeleteAliasByUserID(ctx context.Context, userID string) error {
	err := s.Store.DeleteAliasByUserID(ctx, userID)
	if err != nil {
		log.Printf("error deleting alias: %s", err.Error())
		return ErrDeleteAliasByUserID
	}

	return nil
}
