package service

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
	"ivpn.net/email-service/internal/client/mailer"
	"ivpn.net/email-service/internal/model"
	"ivpn.net/email-service/internal/utils"
)

var (
	ErrGetRecipient      = errors.New("could not get recipient by ID")
	ErrGetRecipients     = errors.New("could not get recipients by user ID")
	ErrPostRecipient     = errors.New("could not create recipient")
	ErrUpdateRecipient   = errors.New("could not update recipient")
	ErrDeleteRecipient   = errors.New("could not delete recipient")
	ErrActivateRecipient = errors.New("could not activate recipient")
)

type RecipienteStore interface {
	GetRecipient(context.Context, string) (model.Recipient, error)
	GetRecipients(context.Context, string) ([]model.Recipient, error)
	PostRecipient(context.Context, model.Recipient) (model.Recipient, error)
	UpdateRecipient(context.Context, model.Recipient) error
	DeleteRecipient(context.Context, string) error
	ActivateRecipient(context.Context, string) error
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
	err := recipient.Validate()
	if err != nil {
		log.Printf("error creating recipient: %s", err.Error())
		return err
	}

	recipient, err = s.Store.PostRecipient(ctx, recipient)
	if err != nil {
		log.Printf("error creating recipient: %s", err.Error())
		switch {
		case errors.Is(err, gorm.ErrDuplicatedKey):
			return model.ErrDuplicateRecipient
		default:
			return ErrPostRecipient
		}
	}

	otp, err := utils.CreateOTP()
	if err != nil {
		log.Printf("error creating recipient: %s", err.Error())
		return ErrCreateOTP
	}

	err = s.Cache.Set(ctx, "activation_recipient_"+recipient.ID, otp.Hash, s.Cfg.Service.OTPExpiration)
	if err != nil {
		log.Printf("error creating recipient: %s", err.Error())
		return ErrSaveOTP
	}

	utils.Background(func() {
		mailer := mailer.New(s.Cfg.SMTPClient)
		err = mailer.Send(recipient.Email, "Activate recipient", otp.Secret)
		if err != nil {
			log.Printf("error creating recipient: %s", err.Error())
		}
	})

	return nil
}

func (s *Service) SendRecipientOTP(ctx context.Context, ID string) error {
	recipient, err := s.GetRecipient(ctx, ID)
	if err != nil {
		log.Printf("error sending OTP: %s", err.Error())
		return ErrGetRecipient
	}

	otp, err := utils.CreateOTP()
	if err != nil {
		log.Printf("error sending OTP: %s", err.Error())
		return ErrCreateOTP
	}

	err = s.Cache.Set(ctx, "activation_recipient_"+ID, otp.Hash, s.Cfg.Service.OTPExpiration)
	if err != nil {
		log.Printf("error sending OTP: %s", err.Error())
		return ErrSaveOTP
	}

	utils.Background(func() {
		mailer := mailer.New(s.Cfg.SMTPClient)
		err = mailer.Send(recipient.Email, "Activate recipient", otp.Secret)
		if err != nil {
			log.Printf("error sending OTP: %s", err.Error())
		}
	})

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

func (s *Service) ActivateRecipient(ctx context.Context, ID string, otp string) error {
	err := utils.ValidateOTP(otp)
	if err != nil {
		log.Printf("error activating recipient: %s", err.Error())
		return err
	}

	hash, err := s.Cache.Get(ctx, "activation_recipient_"+ID)
	if err != nil {
		log.Printf("error activating recipient: %s", err.Error())
		return ErrExpiredOTP
	}

	if hash != utils.HashOTP(otp) {
		return ErrIncorrectOTP
	}

	err = s.Store.ActivateRecipient(ctx, ID)
	if err != nil {
		log.Printf("error activating recipient: %s", err.Error())
		return ErrActivateRecipient
	}

	err = s.Cache.Del(ctx, "activation_recipient_"+ID)
	if err != nil {
		log.Printf("error activating recipient: %s", err.Error())
	}

	return nil
}
