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
	msg, err := parse(from, to, data)
	if err != nil {
		return err
	}

	for _, to := range msg.To {
		recipient, alias, err := s.findRecipient(to)
		if err != nil {
			continue
		}

		mailer := mailer.New(s.Cfg.SMTPClient)
		mailer.Sender = msg.From
		err = mailer.Send(recipient, msg.Subject, msg.Body)
		if err != nil {
			log.Println("error forwarding message", err)
			continue
		}

		err = s.PostMessage(context.Background(), model.Message{
			AliasID: alias.ID,
			UserID:  alias.UserID,
			Type:    model.Forward,
			Size:    len(data),
		})
		if err != nil {
			log.Println("error saving message", err)
		}
	}

	return err
}

func (s *Service) findRecipient(email string) (string, *model.Alias, error) {
	name := email[:strings.Index(email, "@")]
	alias, err := s.GetAliasByName(name)
	if err != nil {
		return "", nil, err
	}

	if !alias.Enabled {
		return "", nil, ErrDisabledAlias
	}

	sub, err := s.GetSubscription(context.Background(), alias.UserID)
	if err != nil {
		return "", nil, err
	}

	if sub.ActiveUntil.Before(time.Now()) {
		return "", nil, ErrInactiveSubscription
	}

	rcps, err := s.GetRecipientsByIDs(context.Background(), alias.Recipients)
	if err != nil || len(rcps) == 0 {
		return "", nil, ErrNoRecipients
	}

	recipient := rcps[0]
	if !recipient.IsActive {
		return "", nil, ErrInactiveRecipient
	}

	return recipient.Email, &alias, nil
}

func parse(from string, to []string, data []byte) (Message, error) {
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
