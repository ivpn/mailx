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

	"ivpn.net/email-service/internal/client/mailer"
	"ivpn.net/email-service/internal/model"
	"ivpn.net/email-service/internal/utils"
)

var (
	ErrInactiveSubscription = errors.New("inactive subscription")
	ErrDisabledAlias        = errors.New("disabled alias")
	ErrNoRecipients         = errors.New("no recipients")
	ErrInactiveRecipient    = errors.New("inactive recipient")
)

type Message struct {
	From    string
	To      []string
	Subject string
	Body    string
}

func (s *Service) ProcessMessage(origin net.Addr, from string, to []string, data []byte) error {
	msg, err := parseMessage(from, to, data)
	if err != nil {
		log.Println("error parsing message", err)
		return err
	}

	for _, to := range msg.To {
		recipient, alias, msgType, err := s.findRecipient(to)
		if err != nil {
			log.Println("error processing message", err)
			continue
		}

		err = s.queueMessage(from, recipient, data, alias, msgType)
		if err != nil {
			log.Println("error queueing message", err)
			continue
		}

		s.saveMessage(alias, msgType, data)
	}

	return err
}

func (s *Service) queueMessage(from string, to string, data []byte, alias model.Alias, msgType model.MessageType) error {
	mailer := mailer.New(s.Cfg.SMTPClient)

	if msgType == model.Forward {
		templateData := map[string]interface{}{
			"alias": alias.Name,
			"from":  from,
		}
		err := mailer.Forward(from, to, data, "header.html", templateData)
		if err != nil {
			log.Println("error forwarding message", err)
			return err
		}
	} else {
		err := mailer.Reply(from, to, data)
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

func (s *Service) findRecipient(email string) (string, model.Alias, model.MessageType, error) {
	name, respondTo := getRespondTo(email)

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

	err = utils.ValidateEmail(respondTo)
	if err == nil {
		return respondTo, alias, model.Reply, nil
	}

	r := strings.Split(alias.Recipients, ",")[0]

	return r, alias, model.Forward, nil
}

func getRespondTo(email string) (string, string) {
	alias := email

	// Get respond to email between "+" and "@"
	rcp := email[strings.Index(email, "+")+1 : strings.Index(email, "@")]

	// Check if respond to email is not empty and contains "="
	if rcp != "" && strings.Contains(rcp, "=") {
		// Replace "=" with "@" to get valid respond to email
		rcp = strings.Replace(rcp, "=", "@", 1)

		// Get alias name up to "+" and domain after "@"
		alias = email[:strings.Index(email, "+")] + email[strings.Index(email, "@"):]
	} else {
		rcp = ""
	}

	return alias, rcp
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

	return Message{
		From:    from,
		To:      to,
		Subject: subject,
		Body:    body,
	}, nil
}
