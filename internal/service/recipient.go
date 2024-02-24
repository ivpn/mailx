package service

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
	"ivpn.net/email-service/internal/model"
)

var (
	ErrGetRecipient    = errors.New("could not get recipient by ID")
	ErrGetRecipients   = errors.New("could not get recipients by user ID")
	ErrPostRecipient   = errors.New("could not post recipient")
	ErrUpdateRecipient = errors.New("could not update recipient")
	ErrDeleteRecipient = errors.New("could not delete recipient")
	ErrVerifyRecipient = errors.New("could not verify recipient")
)

type RecipienteStore interface {
	GetRecipient(context.Context, string) (model.Recipient, error)
	GetRecipients(context.Context, string) ([]model.Recipient, error)
	PostRecipient(context.Context, model.Recipient) error
	UpdateRecipient(context.Context, model.Recipient) error
	DeleteRecipient(context.Context, string) error
	VerifyRecipient(context.Context, string, string) (model.Recipient, error)
}

func (s *Service) GetRecipient(ctx context.Context, ID string) (model.Recipient, error) {
	rcp, err := s.Store.GetRecipient(ctx, ID)
	if err != nil {
		log.Printf("an error occured fetching the recipient: %s", err.Error())
		return model.Recipient{}, ErrGetRecipient
	}

	return rcp, nil
}

func (s *Service) GetRecipients(ctx context.Context, userID string) ([]model.Recipient, error) {
	rcps, err := s.Store.GetRecipients(ctx, userID)
	if err != nil {
		log.Printf("an error occured fetching the recipients: %s", err.Error())
		return []model.Recipient{}, ErrGetRecipients
	}

	return rcps, nil
}

func (s *Service) PostRecipient(ctx context.Context, recipient model.Recipient) error {
	err := s.Store.PostRecipient(ctx, recipient)
	if err != nil {
		log.Printf("an error occurred creating the recipient: %s", err.Error())
		switch {
		case errors.Is(err, gorm.ErrDuplicatedKey):
			return model.ErrDuplicateRecipient
		default:
			return ErrPostRecipient
		}
	}

	return nil
}

func (s *Service) UpdateRecipient(ctx context.Context, recipient model.Recipient) error {
	err := s.Store.UpdateRecipient(ctx, recipient)
	if err != nil {
		log.Printf("an error occurred updating the recipient: %s", err.Error())
		return ErrUpdateRecipient
	}

	return nil
}

func (s *Service) DeleteRecipient(ctx context.Context, ID string) error {
	err := s.Store.DeleteRecipient(ctx, ID)
	if err != nil {
		log.Printf("an error occurred deleting the recipient: %s", err.Error())
		return ErrDeleteRecipient
	}

	return nil
}

func (s *Service) VerifyRecipient(ctx context.Context, ID string, verification string) (model.Recipient, error) {
	rcp, err := s.Store.VerifyRecipient(ctx, ID, verification)
	if err != nil {
		log.Printf("an error occurred verifying the recipient: %s", err.Error())
		return model.Recipient{}, ErrVerifyRecipient
	}

	return rcp, nil
}
