package service

import (
	"context"
	"errors"
	"log"
)

var (
	ErrGetRecipient    = errors.New("could not get recipient by ID")
	ErrPostRecipient   = errors.New("could not post recipient")
	ErrUpdateRecipient = errors.New("could not update recipient")
	ErrDeleteRecipient = errors.New("could not delete recipient")
)

type Recipient struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Verification string `json:"verification"`
}

func (s *Service) GetRecipient(ctx context.Context, ID string) (Recipient, error) {
	rcp, err := s.Store.GetRecipient(ctx, ID)
	if err != nil {
		log.Printf("an error occured fetching the recipient: %s", err.Error())
		return Recipient{}, ErrGetRecipient
	}

	return rcp, nil
}

func (s *Service) PostRecipient(ctx context.Context, rcp Recipient) (Recipient, error) {
	rcp, err := s.Store.PostRecipient(ctx, rcp)
	if err != nil {
		log.Printf("an error occurred creating the recipient: %s", err.Error())
		return Recipient{}, ErrPostRecipient
	}

	return rcp, nil
}

func (s *Service) UpdateRecipient(ctx context.Context, ID string, newRecipient Recipient) (Recipient, error) {
	rcp, err := s.Store.UpdateRecipient(ctx, ID, newRecipient)
	if err != nil {
		log.Printf("an error occurred updating the recipient: %s", err.Error())
		return Recipient{}, ErrUpdateRecipient
	}

	return rcp, nil
}

func (s *Service) DeleteRecipient(ctx context.Context, ID string) error {
	err := s.Store.DeleteRecipient(ctx, ID)
	if err != nil {
		log.Printf("an error occurred deleting the recipient: %s", err.Error())
		return ErrDeleteRecipient
	}

	return nil
}
