package service

import (
	"context"
	"errors"
	"log"

	"ivpn.net/email/api/internal/client/mailer"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
)

var (
	ErrInactiveSubscription = errors.New("Subscription is inactive.")
	ErrDisabledAlias        = errors.New("This alias is disabled.")
	ErrNoRecipients         = errors.New("No recipients found.")
	ErrNoVerifiedRecipients = errors.New("No verified recipients found for sender address.")
	ErrInactiveRecipient    = errors.New("The recipient is inactive.")
)

func (s *Service) ProcessMessage(data []byte) error {
	msg, err := model.ParseMsg(data)
	if err != nil {
		log.Println("error parsing message", err)
		return err
	}

	// Bounce
	if msg.Type == model.FailBounce {
		alias, err := s.FindAlias(msg.From)
		if err != nil {
			log.Println("error processing bounce", err)
			return err
		}

		err = s.ProcessBounce(alias.UserID, alias.ID, data, msg)
		if err != nil {
			log.Println("error processing bounce", err)
			return err
		}

		return nil
	}

	for _, to := range msg.To {
		recipients, alias, relayType, err := s.FindRecipients(msg.From, to, msg.Type)
		if err != nil {
			log.Println("error processing message", err)

			// Handle ErrNoVerifiedRecipients
			if errors.Is(err, ErrNoVerifiedRecipients) {
				settings, err := s.GetSettings(context.Background(), alias.UserID)
				if err != nil {
					log.Println("error getting settings", err)
					continue
				}

				if settings.LogDiscard {
					err := s.ProcessDiscard(alias, msg.From, to, ErrNoVerifiedRecipients.Error())
					if err != nil {
						log.Println("error processing discard", err)
					}
				}
			}

			continue
		}

		sub, err := s.GetSubscription(context.Background(), alias.UserID)
		if err != nil {
			log.Println("error getting subscription", err)
			continue
		}

		// Forward
		if relayType == model.Forward && sub.PendingDelete() {
			log.Println("inactive subscription for forward")
			continue
		}

		// Reply | Send
		if relayType != model.Forward && !sub.ActiveStatus() {
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
				err = s.QueueMessage(msg.From, msg.FromName, settings.FromName, recipient, data, alias, relayType)
				if err != nil {
					log.Println("error queueing message", err)
					return
				}

				err = s.SaveMessage(context.Background(), alias, relayType)
				if err != nil {
					log.Println("error saving message", err)
				}
			})
		}
	}

	return err
}

func (s *Service) QueueMessage(from string, fromName string, settingsFromName string, rcp model.Recipient, data []byte, alias model.Alias, msgType model.MessageType) error {
	mailer := mailer.New(s.Cfg.SMTPClient)

	// Queue Forward
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
		// Queue Reply | Send
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
