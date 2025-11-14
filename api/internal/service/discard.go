package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
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

func (s *Service) ProcessDiscard(alias model.Alias, from string, destination string, message string) error {
	discard := model.Discard{
		ID:          uuid.New().String(),
		CreatedAt:   time.Now(),
		UserID:      alias.UserID,
		AliasID:     alias.ID,
		From:        from,
		Destination: destination,
		Message:     message,
	}

	err := s.PostDiscard(context.Background(), discard)
	if err != nil {
		log.Printf("error processing discard: %s", err.Error())
		return err
	}

	return nil
}
