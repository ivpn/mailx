package service

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
	"ivpn.net/email/api/internal/client/mailer"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
)

var (
	ErrGetRecipient            = errors.New("could not get recipient by ID")
	ErrGetRecipients           = errors.New("could not get recipients by user ID")
	ErrPostRecipient           = errors.New("could not create recipient")
	ErrMaxExceededRecipient    = errors.New("maximum number of recipients reached")
	ErrUpdateRecipient         = errors.New("could not update recipient")
	ErrDeleteRecipient         = errors.New("could not delete recipient")
	ErrDeleteRecipientByUserID = errors.New("could not delete recipient by user ID")
	ErrActivateRecipient       = errors.New("could not activate recipient")
)

type RecipientsStore interface {
	GetRecipient(context.Context, string, string) (model.Recipient, error)
	GetRecipientByEmail(context.Context, string, string) (model.Recipient, error)
	GetRecipientsCountByEmail(context.Context, string) (int, error)
	GetRecipients(context.Context, string) ([]model.Recipient, error)
	GetRecipientsCount(context.Context, string) (int, error)
	GetVerifiedRecipients(context.Context, string, string) ([]model.Recipient, error)
	PostRecipient(context.Context, model.Recipient) (model.Recipient, error)
	UpdateRecipient(context.Context, model.Recipient) error
	DeleteRecipient(context.Context, string, string) error
	ActivateRecipient(context.Context, string, string) error
	DeleteRecipientByUserID(context.Context, string) error
}

func (s *Service) GetRecipient(ctx context.Context, ID string, userID string) (model.Recipient, error) {
	rcp, err := s.Store.GetRecipient(ctx, ID, userID)
	if err != nil {
		log.Printf("an error occured fetching the recipient: %s", err.Error())
		return model.Recipient{}, ErrGetRecipient
	}

	return rcp, nil
}

func (s *Service) GetRecipientByEmail(ctx context.Context, email string, userID string) (model.Recipient, error) {
	rcp, err := s.Store.GetRecipientByEmail(ctx, email, userID)
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

	for i, rcp := range rcps {
		if rcp.PGPKey != "" {
			rcps[i].PGPKey = utils.HashPGPKey(rcp.PGPKey)
		}
	}

	return rcps, nil
}

func (s *Service) GetVerifiedRecipients(ctx context.Context, recipientEmails string, userID string) ([]model.Recipient, error) {
	rcps, err := s.Store.GetVerifiedRecipients(ctx, recipientEmails, userID)
	if err != nil {
		log.Printf("an error occured fetching the recipients: %s", err.Error())
		return []model.Recipient{}, ErrGetRecipients
	}

	return rcps, nil
}

func (s *Service) PostRecipient(ctx context.Context, recipient model.Recipient) error {
	count, err := s.Store.GetRecipientsCount(ctx, recipient.UserID)
	if err != nil {
		log.Printf("error creating recipient: %s", err.Error())
		return ErrPostRecipient
	}

	if count >= s.Cfg.Service.MaxRecipients {
		return ErrMaxExceededRecipient
	}

	recipient, err = s.Store.PostRecipient(ctx, recipient)
	if err != nil {
		log.Printf("error creating recipient: %s", err.Error())
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return model.ErrDuplicateRecipient
		} else {
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
		data := map[string]any{
			"otp":  otp.Secret,
			"from": s.Cfg.SMTPClient.SenderName,
		}
		mailer := mailer.New(s.Cfg.SMTPClient)
		mailer.Sender = s.Cfg.SMTPClient.Sender
		mailer.SenderName = s.Cfg.SMTPClient.SenderName
		err = mailer.SendTemplate(recipient.Email, "["+mailer.SenderName+"] Verify Recipient Notification", "otp_recipient.tmpl", data)
		if err != nil {
			log.Printf("error creating recipient: %s", err.Error())
		}
	})

	return nil
}

func (s *Service) SendRecipientOTP(ctx context.Context, ID string, userID string) error {
	recipient, err := s.GetRecipient(ctx, ID, userID)
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
		data := map[string]any{
			"otp":  otp.Secret,
			"from": s.Cfg.SMTPClient.SenderName,
		}
		mailer := mailer.New(s.Cfg.SMTPClient)
		mailer.Sender = s.Cfg.SMTPClient.Sender
		mailer.SenderName = s.Cfg.SMTPClient.SenderName
		err = mailer.SendTemplate(recipient.Email, "["+mailer.SenderName+"] Verify Recipient Notification", "otp_recipient.tmpl", data)
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

func (s *Service) DeleteRecipient(ctx context.Context, ID string, userID string) error {
	// Get recipient
	recipient, err := s.Store.GetRecipient(ctx, ID, userID)
	if err != nil {
		return err
	}

	// Get aliases
	aliases, err := s.Store.GetAliases(ctx, userID, 0, 0, "", "", "")
	if err != nil {
		return err
	}

	// Delete recipient from each alias
	// Disable alias if no recipients left
	for _, alias := range aliases {
		if strings.Contains(alias.Recipients, recipient.Email) {
			r := alias.Recipients
			r = strings.Replace(r, recipient.Email+",", "", -1)
			r = strings.Replace(r, ","+recipient.Email, "", -1)
			r = strings.Replace(r, recipient.Email, "", -1)
			alias.Recipients = r
			alias.Enabled = r != ""

			// Update alias
			err = s.Store.UpdateAlias(ctx, alias)
			if err != nil {
				return err
			}
		}
	}

	err = s.Store.DeleteRecipient(ctx, ID, userID)
	if err != nil {
		log.Printf("an error occurred deleting the recipient: %s", err.Error())
		return ErrDeleteRecipient
	}

	return nil
}

func (s *Service) ActivateRecipient(ctx context.Context, ID string, userID string, otp string) error {
	hash, err := s.Cache.Get(ctx, "activation_recipient_"+ID)
	if err != nil {
		log.Printf("error activating recipient: %s", err.Error())
		return ErrExpiredOTP
	}

	if !utils.MatchOTP(otp, hash) {
		return ErrIncorrectOTP
	}

	err = s.Store.ActivateRecipient(ctx, ID, userID)
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

func (s *Service) DeleteRecipientByUserID(ctx context.Context, userID string) error {
	err := s.Store.DeleteRecipientByUserID(ctx, userID)
	if err != nil {
		log.Printf("an error occurred deleting the recipient: %s", err.Error())
		return ErrDeleteRecipientByUserID
	}

	return nil
}
