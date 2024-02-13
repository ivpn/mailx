package service

import (
	"context"
	"errors"
	"log"

	"ivpn.net/email-service/internal/model"
)

var (
	ErrGetAlias       = errors.New("could not get alias by ID")
	ErrGetAliases     = errors.New("could not get aliass by recipient ID")
	ErrGetAliasByName = errors.New("could not get alias by name")
	ErrPostAlias      = errors.New("could not post alias")
	ErrUpdateAlias    = errors.New("could not update alias")
	ErrDeleteAlias    = errors.New("could not delete alias")
)

type AliasStore interface {
	GetAlias(context.Context, string) (model.Alias, error)
	GetAliases(context.Context, string) ([]model.Alias, error)
	GetAliasByName(string) (model.Alias, error)
	PostAlias(context.Context, model.Alias) error
	UpdateAlias(context.Context, string, model.Alias) error
	DeleteAlias(context.Context, string) error
}

func (s *Service) GetAlias(ctx context.Context, ID string) (model.Alias, error) {
	alias, err := s.Store.GetAlias(ctx, ID)
	if err != nil {
		log.Printf("an error occured fetching the alias: %s", err.Error())
		return model.Alias{}, ErrGetAlias
	}

	return alias, nil
}

func (s *Service) GetAliases(ctx context.Context, recipientID string) ([]model.Alias, error) {
	aliases, err := s.Store.GetAliases(ctx, recipientID)
	if err != nil {
		log.Printf("an error occured fetching the aliass: %s", err.Error())
		return []model.Alias{}, ErrGetAliases
	}

	return aliases, nil
}

func (s *Service) GetAliasByName(name string) (model.Alias, error) {
	alias, err := s.Store.GetAliasByName(name)
	if err != nil {
		log.Printf("an error occured fetching the alias: %s", err.Error())
		return model.Alias{}, ErrGetAliasByName
	}

	return alias, nil
}

func (s *Service) PostAlias(ctx context.Context, alias model.Alias) error {
	err := s.Store.PostAlias(ctx, alias)
	if err != nil {
		log.Printf("an error occurred creating the alias: %s", err.Error())
		return ErrPostAlias
	}

	return nil
}

func (s *Service) UpdateAlias(ctx context.Context, ID string, newAlias model.Alias) error {
	err := s.Store.UpdateAlias(ctx, ID, newAlias)
	if err != nil {
		log.Printf("an error occurred updating the alias: %s", err.Error())
		return ErrUpdateAlias
	}

	return nil
}

func (s *Service) DeleteAlias(ctx context.Context, ID string) error {
	err := s.Store.DeleteAlias(ctx, ID)
	if err != nil {
		log.Printf("an error occurred deleting the alias: %s", err.Error())
		return ErrDeleteAlias
	}

	return nil
}
