package service

import (
	"context"
	"errors"
	"log"
	"strings"

	"ivpn.net/email/api/internal/client/mailer"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
)

var (
	ErrInactiveSubscription = errors.New("Subscription is inactive.")
	ErrDisabledAlias        = errors.New("This alias is disabled.")
	ErrNoRecipients         = errors.New("No verified recipients available.")
	ErrInactiveRecipient    = errors.New("The recipient is inactive.")
)

func (s *Service) ProcessMessage(data []byte) error {
	msg, err := model.ParseMsg(data)
	if err != nil {
		log.Println("error parsing message", err)
		return err
	}

	for _, to := range msg.To {
		recipients, alias, relayType, err := s.findRecipients(msg.From, to, msg.Type)
		if err != nil {
			log.Println("error processing message", err)
			continue
		}

		sub, err := s.GetSubscription(context.Background(), alias.UserID)
		if err != nil {
			log.Println("error getting subscription", err)
			continue
		}

		// Forward
		if relayType == model.Forward && !sub.IsActiveWithGracePeriod(s.Cfg.Service.ForwardGracePeriodDays) {
			log.Println("inactive subscription for forward")
			continue
		}

		// Reply | Send
		if relayType != model.Forward && !sub.IsActive() {
			log.Println("inactive subscription for reply/send")
			continue
		}

		settings, err := s.GetSettings(context.Background(), alias.UserID)
		if err != nil {
			log.Println("error getting settings", err)
			continue
		}

		for _, recipient := range recipients {
			utils.Background(func() {
				err = s.queueMessage(msg.From, msg.FromName, settings.FromName, recipient, data, alias, relayType)
				if err != nil {
					log.Println("error queueing message", err)
					return
				}

				s.saveMessage(alias, relayType, data)
			})
		}
	}

	return err
}

func (s *Service) queueMessage(from string, fromName string, settingsFromName string, rcp model.Recipient, data []byte, alias model.Alias, msgType model.MessageType) error {
	mailer := mailer.New(s.Cfg.SMTPClient)

	// Forward
	if msgType == model.Forward {
		templateData := map[string]any{
			"alias": alias.Name,
			"from":  from,
		}
		generatedFrom := model.GenerateReplyTo(alias.Name, from)
		err := mailer.Forward(generatedFrom, fromName, rcp, data, "header.tmpl", templateData)
		if err != nil {
			log.Println("error forwarding message", err)
			return err
		}
	} else {
		// Reply | Send
		err := s.ValidateSendReplyDailyCount(context.Background(), alias.UserID)
		if err != nil {
			log.Println("error validating send/reply daily count", err)
			return err
		}

		name := alias.FromName
		if name == "" {
			name = settingsFromName
		}

		err = mailer.Reply(alias.Name, name, rcp, data)
		if err != nil {
			log.Println("error sending message", err)
			return err
		}
	}

	return nil
}

func (s *Service) saveMessage(alias model.Alias, msgType model.MessageType, data []byte) {
	err := s.PostMessage(context.Background(), model.Message{
		AliasID: alias.ID,
		UserID:  alias.UserID,
		Type:    msgType,
	})
	if err != nil {
		log.Println("error saving message", err)
	}
}

func (s *Service) findRecipients(from string, email string, msgType model.MessageType) ([]model.Recipient, model.Alias, model.MessageType, error) {
	name, replyTo := model.ParseReplyTo(email)

	alias, err := s.GetAliasByName(name)
	if err != nil {
		return []model.Recipient{}, model.Alias{}, 0, err
	}

	if !alias.Enabled {
		// Block
		s.saveMessage(alias, model.Block, []byte{})
		return []model.Recipient{}, model.Alias{}, 0, ErrDisabledAlias
	}

	err = utils.ValidateEmail(replyTo)
	if err == nil {
		rcps, err := s.GetVerifiedRecipients(context.Background(), from, alias.UserID)
		if err != nil || len(rcps) == 0 {
			return []model.Recipient{}, model.Alias{}, 0, ErrNoRecipients
		}

		return []model.Recipient{{Email: replyTo}}, alias, model.MessageType(msgType), nil
	}

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
