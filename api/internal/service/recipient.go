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
	ErrGetRecipient            = errors.New("Unable to retrieve recipient by ID.")
	ErrGetRecipients           = errors.New("Unable to retrieve recipients by user ID.")
	ErrPostRecipient           = errors.New("Unable to create recipient.")
	ErrMaxExceededRecipient    = errors.New("Maximum number of allowed recipients reached.")
	ErrUpdateRecipient         = errors.New("Unable to update recipient.")
	ErrDeleteRecipient         = errors.New("Unable to delete recipient.")
	ErrDeleteRecipientByUserID = errors.New("Unable to delete recipient for this user.")
	ErrActivateRecipient       = errors.New("Unable to activate recipient.")
)

type RecipientsStore interface {
	GetRecipient(context.Context, string, string) (model.Recipient, error)
	GetRecipientByEmail(context.Context, string, string) (model.Recipient, error)
	CheckDuplicateRecipient(context.Context, string) (bool, error)
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
	sub, err := s.GetSubscription(context.Background(), recipient.UserID)
	if err != nil {
		log.Printf("error fetching subscription: %s", err.Error())
		return ErrPostRecipient
	}

	if !sub.ActiveStatus() {
		log.Println("error creating recipient: subscription is not active")
		return ErrPostRecipient
	}

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
		err = mailer.SendTemplate(recipient.Email, "Verify Your Recipient Email Address", "otp_recipient.tmpl", data)
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
		err = mailer.SendTemplate(recipient.Email, "Verify Your Recipient Email Address", "otp_recipient.tmpl", data)
		if err != nil {
			log.Printf("error sending OTP: %s", err.Error())
		}
	})

	return nil
}

func (s *Service) UpdateRecipient(ctx context.Context, recipient model.Recipient) error {
	sub, err := s.GetSubscription(context.Background(), recipient.UserID)
	if err != nil {
		log.Printf("error fetching subscription: %s", err.Error())
		return ErrUpdateRecipient
	}

	if !sub.ActiveStatus() {
		log.Println("error updating recipient: subscription is not active")
		return ErrUpdateRecipient
	}

	err = s.Store.UpdateRecipient(ctx, recipient)
	if err != nil {
		log.Printf("error updating recipient: %s", err.Error())
		return ErrUpdateRecipient
	}

	return nil
}

func (s *Service) DeleteRecipient(ctx context.Context, ID string, userID string, newRecipients string) error {
	// Get recipient
	recipient, err := s.Store.GetRecipient(ctx, ID, userID)
	if err != nil {
		log.Printf("error deleting recipient, GetRecipient: %s", err.Error())
		return ErrDeleteRecipient
	}

	// Get aliases
	aliases, err := s.Store.GetAliases(ctx, userID, 0, 0, "created_at", "DESC", "", "")
	if err != nil {
		log.Printf("error deleting recipient, GetAliases: %s", err.Error())
		return ErrDeleteRecipient
	}

	// Delete recipient from each alias
	// Disable alias if no recipients left
	for _, alias := range aliases {
		if strings.Contains(alias.Recipients, recipient.Email) {
			r := alias.Recipients
			r = strings.Replace(r, recipient.Email+",", "", -1)
			r = strings.Replace(r, ","+recipient.Email, "", -1)
			r = strings.Replace(r, recipient.Email, "", -1)
			alias.Recipients = model.MergeCommaSeparatedEmails(r, newRecipients)
			alias.Enabled = alias.Recipients != ""

			// Update alias
			err = s.Store.UpdateAlias(ctx, alias)
			if err != nil {
				log.Printf("error deleting recipient, UpdateAlias: %s", err.Error())
				return ErrDeleteRecipient
			}
		}
	}

	err = s.Store.DeleteRecipient(ctx, ID, userID)
	if err != nil {
		log.Printf("error deleting recipient, DeleteRecipient: %s", err.Error())
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

	idLimiter := utils.IDLimiter{
		ID:    ID,
		Label: "recipient_fails",
		Max:   s.Cfg.Service.IdLimiterMax,
		Exp:   s.Cfg.Service.IdLimiterExpiration,
		Cache: s.Cache,
	}

	if !utils.MatchOTP(otp, hash) {
		err = idLimiter.Tick()
		if err != nil {
			log.Printf("error activating recipient: %s", err.Error())
		}

		return ErrIncorrectOTP
	}

	if !idLimiter.IsAllowed() {
		log.Printf("error activating recipient: too many failed attempts")
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

func (s *Service) FindRecipients(from string, to string, msgType model.MessageType) ([]model.Recipient, model.Alias, model.MessageType, error) {
	// Extract alias name from the "to" email
	name, replyTo := model.ParseReplyTo(to)
	alias, err := s.GetAliasByName(name)
	if err != nil {
		return []model.Recipient{}, model.Alias{}, 0, err
	}

	// Handle disabled alias
	if !alias.Enabled {
		err = s.SaveMessage(context.Background(), alias, model.Block)
		if err != nil {
			log.Println("error saving message", err)
		}

		return []model.Recipient{}, model.Alias{}, 0, ErrDisabledAlias
	}

	// Handle Reply | Send
	err = utils.ValidateEmail(replyTo)
	if err == nil {
		rcps, err := s.GetVerifiedRecipients(context.Background(), from, alias.UserID)
		if err != nil || len(rcps) == 0 {
			return []model.Recipient{}, alias, 0, ErrNoVerifiedRecipients
		}

		return []model.Recipient{{Email: replyTo}}, alias, model.MessageType(msgType), nil
	}

	// Handle Forward
	rcps, err := s.GetRecipients(context.Background(), alias.UserID)
	if err != nil || len(rcps) == 0 {
		return []model.Recipient{}, model.Alias{}, 0, ErrNoRecipients
	}

	var recipients []model.Recipient
	for _, rcp := range rcps {
		if strings.Contains(alias.Recipients, rcp.Email) {
			recipients = append(recipients, rcp)
		}
	}

	return recipients, alias, model.Forward, nil
}
