package service

import (
	"context"
	"errors"
	"log"

	"golang.org/x/sync/errgroup"
	"ivpn.net/email/api/internal/client/mailer"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
)

var (
	ErrInactiveSubscription = errors.New("Subscription is inactive.")
	ErrNoRecipients         = errors.New("No recipients found.")
	ErrNoVerifiedRecipients = errors.New("Sender is not a verified address.")
	ErrInactiveRecipient    = errors.New("The recipient is inactive.")
)

func (s *Service) ProcessMessage(data []byte) error {
	msg, err := model.ParseMsg(data)
	if err != nil {
		log.Println("error parsing message:", err)
		return err
	}

	// Bounce
	if msg.Type == model.FailBounce {
		alias, err := s.FindAlias(msg.From)
		if err != nil {
			log.Println("error processing bounce:", err, alias.Name)
			return err
		}

		err = s.ProcessBounceLog(alias.UserID, alias.ID, data, msg)
		if err != nil {
			log.Println("error processing bounce:", err, alias.Name)
			return err
		}

		return nil
	}

	// Verify Email Authentication
	pass, err := utils.VerifyEmailAuth(data)
	if err != nil {
		log.Println("email authentication failed:", err)
	}
	if !pass {
		// Fail silently so unauthenticated emails are not kept in postfix queue
		return nil
	}

	var g errgroup.Group

	for _, to := range msg.To {
		recipients, alias, relayType, err := s.FindRecipients(msg.From, to, msg.Type)
		if err != nil {
			log.Println("error processing message:", err, alias.Name)

			// Handle ErrNoVerifiedRecipients
			if errors.Is(err, ErrNoVerifiedRecipients) {
				settings, err := s.GetSettings(context.Background(), alias.UserID)
				if err != nil {
					log.Println("error getting settings", err)
					continue
				}

				if settings.LogIssues {
					err := s.ProcessDiscardLog(alias, msg.From, to, ErrNoVerifiedRecipients.Error(), model.UnauthorisedSend)
					if err != nil {
						log.Println("error processing discard log", err)
					}
				}
			}

			// Handle ErrDisabledAlias
			if errors.Is(err, ErrDisabledAlias) {
				settings, err := s.GetSettings(context.Background(), alias.UserID)
				if err != nil {
					log.Println("error getting settings", err)
					continue
				}

				if settings.LogIssues {
					err := s.ProcessDiscardLog(alias, msg.From, to, ErrDisabledAlias.Error(), model.DisabledAlias)
					if err != nil {
						log.Println("error processing discard log", err)
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
			g.Go(func() error {
				err = s.QueueMessage(msg.From, msg.FromName, recipient, data, alias, relayType, settings)
				if err != nil {
					return err
				}

				if err := s.SaveMessage(context.Background(), alias, relayType); err != nil {
					log.Println("error saving message", err)
				}

				return nil
			})
		}
	}

	// Wait for all goroutines and return first error (if any)
	return g.Wait()
}

func (s *Service) QueueMessage(from string, fromName string, rcp model.Recipient, data []byte, alias model.Alias, msgType model.MessageType, settings model.Settings) error {
	mailer := mailer.New(s.Cfg.SMTPClient)

	// Queue Forward
	if msgType == model.Forward {
		templateData := map[string]any{
			"alias": alias.Name,
			"from":  from,
		}
		generatedFrom := model.GenerateReplyTo(alias.Name, from)
		err := mailer.Forward(generatedFrom, fromName, rcp, data, "header.tmpl", templateData, settings)
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
			name = settings.FromName
		}

		err = mailer.Reply(alias.Name, name, rcp, data)
		if err != nil {
			log.Println("error sending message", err)
			return err
		}
	}

	return nil
}
