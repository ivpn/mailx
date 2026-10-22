package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"time"

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

	// queueMessage hands an email to the mailer for one recipient; tests replace it.
	queueMessage = (*Service).QueueMessage
	// noRecipientsReason is the diagnostic shown when ErrNoRecipients drops a message, spelling
	// out where a recipient can be set since the bare error doesn't say.
	noRecipientsReason = "No verified recipient to forward to. Set a verified recipient on the alias, the domain's catch-all recipient, or the default recipient in Settings."
	// disabledAliasReason and disabledDomainReason are the diagnostics shown when ErrDisabledAlias
	// and ErrDisabledDomain drop a message; the bare errors are meant to be followed by a name.
	disabledAliasReason  = "The alias is disabled. Enable it to receive emails again."
	disabledDomainReason = "The alias's custom domain is disabled or not verified. Enable or verify the domain to receive emails again."
)

func (s *Service) ProcessMessage(data []byte, envelopeRecipient string) error {
	msg, parseErr := model.ParseMsg(data)
	if envelopeRecipient == "" {
		// Caller didn't pass the actual envelope recipient (e.g. Postfix pipe not
		// yet updated to pass it); fall back to guessing it from headers.
		envelopeRecipient = msg.EnvelopeRecipient
	}

	if parseErr != nil {
		if errors.Is(parseErr, model.ErrExtractOriginalFrom) {
			// Fail silently so bounce messages are not kept in postfix queue
			return nil
		}

		for _, to := range utils.SelectTargets(msg.To, envelopeRecipient) {
			_, alias, _, err := s.FindRecipients(msg.From, to, msg.Type)
			if alias.UserID == "" {
				continue
			}

			settings, err := s.GetSettings(context.Background(), alias.UserID)
			if err != nil {
				log.Println("error getting settings", err)
				continue
			}

			if settings.LogIssues {
				err := s.ProcessDiagnosticLog(alias, msg.From, to, parseErr.Error(), model.DeferredDelivery)
				if err != nil {
					log.Println("error processing diagnostic log", err)
				}
			}

			log.Println("error parsing message [alias:", alias.ID, "]:", parseErr)
		}

		// Disable fallback logging to avoid cluttering
		// log.Println("error parsing message [data]:", parseErr)
		return parseErr
	}

	// Bounce
	if msg.Type == model.FailBounce {
		alias, err := s.FindAlias(msg.From)
		if err != nil {
			log.Println("error processing bounce:", err, alias.ID)
			// Fail silently so bounce messages are not kept in postfix queue
			return nil
		}

		err = s.ProcessBounceLog(alias.UserID, alias.ID, data, msg)
		if err != nil {
			log.Println("error processing bounce:", err, alias.ID)
		}

		// Fail silently so bounce messages are not kept in postfix queue
		return nil
	}

	// Verify Email Authentication
	pass, err := utils.VerifyEmailAuth(data, utils.EmailAuthConfig{
		TrustedRelayDomains: s.Cfg.API.TrustedRelayDomains,
		DisableAlignment:    s.Cfg.API.DisableEmailAuthAlignment,
		LogMismatch:         s.Cfg.API.LogEmailAuthMismatch,
	})
	if err != nil {
		log.Println("email authentication failed:", err)
	}
	if !pass {
		// Fail silently so unauthenticated emails are not kept in postfix queue
		return nil
	}

	var g errgroup.Group

	for _, to := range utils.SelectTargets(msg.To, envelopeRecipient) {
		recipients, alias, relayType, err := s.FindRecipients(msg.From, to, msg.Type)
		// No recipient means the email would be dropped silently, so log it as ErrNoRecipients
		if err == nil && len(recipients) == 0 {
			err = ErrNoRecipients
		}
		if err != nil {
			log.Println("error processing message:", err, alias.ID)

			// Handle ErrNoVerifiedRecipients
			if errors.Is(err, ErrNoVerifiedRecipients) {
				settings, err := s.GetSettings(context.Background(), alias.UserID)
				if err != nil {
					log.Println("error getting settings", err)
					continue
				}

				if settings.LogIssues {
					err := s.ProcessDiagnosticLog(alias, msg.From, to, ErrNoVerifiedRecipients.Error(), model.UnauthorisedSend)
					if err != nil {
						log.Println("error processing diagnostic log", err)
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
					err := s.ProcessDiagnosticLog(alias, msg.From, to, disabledAliasReason, model.DisabledAlias)
					if err != nil {
						log.Println("error processing diagnostic log", err)
					}
				}
			}

			// Handle ErrDisabledDomain
			if errors.Is(err, ErrDisabledDomain) {
				settings, err := s.GetSettings(context.Background(), alias.UserID)
				if err != nil {
					log.Println("error getting settings", err)
					continue
				}

				if settings.LogIssues {
					err := s.ProcessDiagnosticLog(alias, msg.From, to, disabledDomainReason, model.DisabledDomain)
					if err != nil {
						log.Println("error processing diagnostic log", err)
					}
				}
			}

			// Handle ErrNoRecipients
			if errors.Is(err, ErrNoRecipients) {
				settings, err := s.GetSettings(context.Background(), alias.UserID)
				if err != nil {
					log.Println("error getting settings", err)
					continue
				}

				if settings.LogIssues {
					err := s.ProcessDiagnosticLog(alias, msg.From, to, noRecipientsReason, model.NoRecipients)
					if err != nil {
						log.Println("error processing diagnostic log", err)
					}
				}
			}

			continue
		}

		// Paused alias - forwards are discarded until paused_until. Replies/sends still work.
		if relayType == model.Forward && alias.IsPaused(time.Now().UTC()) {
			s.DiscardPausedMessage(alias, msg.From, to)
			continue
		}

		// Blocklist - checked before the inbound alias below so a blocked sender can't provision
		// catch-all aliases. Replies/sends come from the user's own recipients, so only forwards apply.
		if relayType == model.Forward {
			if item, blocked := s.FindBlocklistMatch(context.Background(), alias.UserID, msg.From); blocked {
				s.BlockMessage(alias, msg.From, to, item)
				continue
			}
		}

		sub, err := s.GetSubscription(context.Background(), alias.UserID)
		if err != nil {
			log.Println("error getting subscription", err)
			continue
		}

		// Reply | Send
		if relayType != model.Forward && !sub.ActiveStatus() {
			settings, err := s.GetSettings(context.Background(), alias.UserID)
			if err != nil {
				log.Println("error getting settings", err)
				continue
			}

			if settings.LogIssues {
				err := s.ProcessDiagnosticLog(alias, msg.From, to, ErrInactiveSubscription.Error(), model.InactiveSubscription)
				if err != nil {
					log.Println("error processing diagnostic log", err)
				}
			}

			continue
		}

		settings, err := s.GetSettings(context.Background(), alias.UserID)
		if err != nil {
			log.Println("error getting settings", err)
			continue
		}

		// Rules - matched against the original email before the inbound alias below, so a rule
		// that blocks the email also stops the alias from being created.
		outcome := s.ApplyRules(context.Background(), alias, msg, relayType, to, data)

		// Handle Inbound Alias
		if alias.Origin == model.Inbound && alias.ID == "" && !outcome.Block {
			inboundAlias, err := s.PostInboundAlias(context.Background(), alias)
			if err == nil {
				alias.BaseModel = inboundAlias.BaseModel
			} else {
				log.Println("error auto-creating inbound alias:", err)
			}
		}

		if outcome.Matched() {
			s.ApplyRuleSideEffects(context.Background(), alias, msg.From, outcome)
		}

		if outcome.Block {
			s.BlockMessageByRule(alias, msg.From, to, relayType, outcome)
			continue
		}

		settings = outcome.Settings(settings)
		recipients = s.RuleRecipients(context.Background(), alias, recipients, outcome)

		toHeader := to
		listUnsubscribe := model.ListUnsubscribeHeaders{}
		if relayType == model.Forward {
			toHeader = utils.CombineForwardTo(to, msg.To, model.GenerateReplyTo)
			listUnsubscribe = s.ListUnsubscribeHeaders(settings, alias, msg)
		}

		for _, recipient := range recipients {
			g.Go(func() (err error) {
				// The API's recover middleware can't catch a panic in this goroutine, and it would
				// stop the server. As an error, Postfix retries the message instead.
				defer func() {
					if r := recover(); r != nil {
						log.Printf("panic queueing message [alias: %s]: %v\n%s", alias.ID, r, debug.Stack())
						err = fmt.Errorf("panic queueing message: %v", r)
					}
				}()

				// Queue Message
				err = queueMessage(s, msg.From, msg.FromName, to, toHeader, recipient, data, alias, relayType, settings, listUnsubscribe, outcome)
				if err != nil {
					return err
				}

				// Save Message for stats
				go func() {
					err := s.SaveMessage(context.Background(), alias, relayType)
					if err != nil {
						log.Println("error saving message", err)
					}
				}()

				return nil
			})
		}
	}

	// Wait for all goroutines and return first error (if any)
	return g.Wait()
}

func (s *Service) QueueMessage(from string, fromName string, to string, toHeader string, rcp model.Recipient, data []byte, alias model.Alias, msgType model.MessageType, settings model.Settings, listUnsubscribe model.ListUnsubscribeHeaders, outcome model.RuleOutcome) error {
	client := mailer.New(s.Cfg.SMTPClient)

	// Queue Forward
	if msgType == model.Forward {
		templateData := map[string]any{
			"alias": to,
			"from":  from,
		}
		if outcome.FromName != "" {
			fromName = outcome.FromName
		}
		generatedFrom := model.GenerateReplyTo(alias.Name, from)
		err := client.Forward(generatedFrom, fromName, toHeader, rcp, data, "header.tmpl", templateData, settings, alias, listUnsubscribe, outcome)
		if err != nil {
			if settings.LogIssues {
				err := s.ProcessDiagnosticLog(alias, from, rcp.Email, err.Error(), model.DeferredDelivery)
				if err != nil {
					log.Println("error processing diagnostic log", err)
				}
			}

			log.Println("error forwarding message [alias:", alias.ID, "]:", mailer.ErrorClass(err))
			return err
		}
	} else {
		// Queue Reply | Send
		err := s.ValidateSendReplyDailyCount(context.Background(), alias.UserID)
		if err != nil {
			log.Println("error validating send/reply daily count", err)
			return err
		}

		err = client.Reply(alias.Name, replyFromName(alias, settings, outcome), rcp, data, alias, outcome)
		if err != nil {
			if settings.LogIssues {
				err := s.ProcessDiagnosticLog(alias, from, rcp.Email, err.Error(), model.DeferredDelivery)
				if err != nil {
					log.Println("error processing diagnostic log", err)
				}
			}

			log.Println("error sending message [alias:", alias.ID, "]:", mailer.ErrorClass(err))
			return err
		}
	}

	return nil
}

// replyFromName is the name a reply or send is shown to come from: a rule's sender name, else
// the alias's, else the user's default.
func replyFromName(alias model.Alias, settings model.Settings, outcome model.RuleOutcome) string {
	if outcome.FromName != "" {
		return outcome.FromName
	}
	if alias.FromName != "" {
		return alias.FromName
	}

	return settings.FromName
}
