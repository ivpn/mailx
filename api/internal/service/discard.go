package service

import (
	"context"
	"errors"
	"log"

	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetDiscards    = errors.New("Unable to retrieve discards for this user.")
	ErrPostDiscard    = errors.New("Unable to create discard.")
	ErrDeleteDiscards = errors.New("Unable to delete discards for this user.")
)

type DiscardStore interface {
	GetDiscards(context.Context, string) ([]model.Discard, error)
	PostDiscard(context.Context, model.Discard) error
	DeleteDiscards(context.Context, string) error
}

func (s *Service) GetDiscards(ctx context.Context, userID string) ([]model.Discard, error) {
	discards, err := s.Store.GetDiscards(ctx, userID)
	if err != nil {
		log.Printf("error getting discards by user ID: %s", err.Error())
		return nil, ErrGetDiscards
	}

	return discards, nil
}

func (s *Service) PostDiscard(ctx context.Context, discard model.Discard) error {
	err := s.Store.PostDiscard(ctx, discard)
	if err != nil {
		log.Printf("error posting discard: %s", err.Error())
		return ErrPostDiscard
	}

	return nil
}

func (s *Service) DeleteDiscards(ctx context.Context, userID string) error {
	err := s.Store.DeleteDiscards(ctx, userID)
	if err != nil {
		log.Printf("error deleting discards by user ID: %s", err.Error())
		return ErrDeleteDiscards
	}

	return nil
}
