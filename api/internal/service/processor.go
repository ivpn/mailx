package service

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net"
	"net/mail"
	"strings"
	"time"

	"ivpn.net/email/api/internal/client/mailer"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
)

var (
	ErrInactiveSubscription = errors.New("inactive subscription")
	ErrDisabledAlias        = errors.New("disabled alias")
	ErrNoRecipients         = errors.New("no recipients")
	ErrInactiveRecipient    = errors.New("inactive recipient")
)

type Message struct {
	From      string
	To        []string
	Subject   string
	Body      string
	RelayType model.MessageType
}

func (s *Service) ProcessMessage(origin net.Addr, from string, to []string, data []byte) error {
	msg, err := parseMessage(from, to, data)
	if err != nil {
		log.Println("error parsing message", err)
		return err
	}

	for _, to := range msg.To {
		recipient, alias, msgType, err := s.findRecipient(to, msg.RelayType)
		if err != nil {
			log.Println("error processing message", err)
			continue
		}

		settings, err := s.GetSettings(context.Background(), alias.UserID)
		if err != nil {
			log.Println("error getting settings", err)
			continue
		}

		utils.Background(func() {
			err = s.queueMessage(from, settings.FromName, recipient, data, alias, msgType)
			if err != nil {
				log.Println("error queueing message", err)
				return
			}

			s.saveMessage(alias, msgType, data)
		})
	}

	return err
}

func (s *Service) queueMessage(from string, settingsFromName string, to string, data []byte, alias model.Alias, msgType model.MessageType) error {
	mailer := mailer.New(s.Cfg.SMTPClient)

	name := alias.FromName
	if name == "" {
		name = settingsFromName
	}

	if msgType == model.Forward {
		templateData := map[string]interface{}{
			"alias": alias.Name,
			"from":  from,
		}
		generatedFrom := model.GenerateReplyTo(alias.Name, from)
		err := mailer.Forward(generatedFrom, name, to, data, "header.tmpl", templateData)
		if err != nil {
			log.Println("error forwarding message", err)
			return err
		}
	} else {
		err := mailer.Reply(alias.Name, name, to, data)
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
		Size:    len(data),
	})
	if err != nil {
		log.Println("error saving message", err)
	}
}

func (s *Service) findRecipient(email string, relayType model.MessageType) (string, model.Alias, model.MessageType, error) {
	name, replyTo := model.ParseReplyTo(email)

	alias, err := s.GetAliasByName(name)
	if err != nil {
		return "", model.Alias{}, 0, err
	}

	if !alias.Enabled {
		s.saveMessage(alias, model.Block, []byte{})
		return "", model.Alias{}, 0, ErrDisabledAlias
	}

	sub, err := s.GetSubscription(context.Background(), alias.UserID)
	if err != nil {
		return "", model.Alias{}, 0, err
	}

	if sub.ActiveUntil.Before(time.Now()) {
		return "", model.Alias{}, 0, ErrInactiveSubscription
	}

	err = utils.ValidateEmail(replyTo)
	if err == nil {
		msgType := model.Send
		if relayType == model.Reply {
			msgType = model.Reply
		}

		return replyTo, alias, msgType, nil
	}

	r := strings.Split(alias.Recipients, ",")[0]

	return r, alias, relayType, nil
}

func parseMessage(from string, to []string, data []byte) (Message, error) {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return Message{}, err
	}

	subject := msg.Header.Get("Subject")
	buf := new(bytes.Buffer)
	buf.ReadFrom(msg.Body)
	body := buf.String()
	relayType := model.Forward

	if isReply(msg) {
		relayType = model.Reply
	}

	return Message{
		From:      from,
		To:        to,
		Subject:   subject,
		Body:      body,
		RelayType: relayType,
	}, nil
}

func isReply(m *mail.Message) bool {
	if m.Header.Get("In-Reply-To") != "" || m.Header.Get("References") != "" {
		return true
	}

	return false
}
