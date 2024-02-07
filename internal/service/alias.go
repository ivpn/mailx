package service

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
)

var (
	ErrGetAlias    = errors.New("could not get alias by ID")
	ErrGetAliases  = errors.New("could not get aliass by recipient ID")
	ErrPostAlias   = errors.New("could not post alias")
	ErrUpdateAlias = errors.New("could not update alias")
	ErrDeleteAlias = errors.New("could not delete alias")
)

type Alias struct {
	gorm.Model
	ID          string `json:"id"`
	RecipientID string `json:"recipient_id"`
	Slug        string `json:"slug"`
}

type AliasStore interface {
	GetAlias(context.Context, string) (Alias, error)
	GetAliases(context.Context, string) ([]Alias, error)
	PostAlias(context.Context, Alias) (Alias, error)
	UpdateAlias(context.Context, string, Alias) (Alias, error)
	DeleteAlias(context.Context, string) error
}

func (s *Service) GetAlias(ctx context.Context, ID string) (Alias, error) {
	rcp, err := s.Store.GetAlias(ctx, ID)
	if err != nil {
		log.Printf("an error occured fetching the alias: %s", err.Error())
		return Alias{}, ErrGetAlias
	}

	return rcp, nil
}

func (s *Service) GetAliass(ctx context.Context, recipientID string) ([]Alias, error) {
	rcps, err := s.Store.GetAliases(ctx, recipientID)
	if err != nil {
		log.Printf("an error occured fetching the aliass: %s", err.Error())
		return []Alias{}, ErrGetAliases
	}

	return rcps, nil
}

func (s *Service) PostAlias(ctx context.Context, rcp Alias) (Alias, error) {
	rcp, err := s.Store.PostAlias(ctx, rcp)
	if err != nil {
		log.Printf("an error occurred creating the alias: %s", err.Error())
		return Alias{}, ErrPostAlias
	}

	return rcp, nil
}

func (s *Service) UpdateAlias(ctx context.Context, ID string, newAlias Alias) (Alias, error) {
	rcp, err := s.Store.UpdateAlias(ctx, ID, newAlias)
	if err != nil {
		log.Printf("an error occurred updating the alias: %s", err.Error())
		return Alias{}, ErrUpdateAlias
	}

	return rcp, nil
}

func (s *Service) DeleteAlias(ctx context.Context, ID string) error {
	err := s.Store.DeleteAlias(ctx, ID)
	if err != nil {
		log.Printf("an error occurred deleting the alias: %s", err.Error())
		return ErrDeleteAlias
	}

	return nil
}

func generateSlug() string {
	rand.Seed(time.Now().UnixNano())

	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 8)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}
