package service

import (
	"bytes"
	"context"
	"net"
	"net/mail"
	"strings"

	"ivpn.net/email-service/internal/client/email"
)

type Message struct {
	From    string
	To      []string
	Subject string
	Body    string
}

func (s *Service) ProcessMessage(origin net.Addr, from string, to []string, data []byte) error {
	msg, err := parse(origin, from, to, data)
	if err != nil {
		return err
	}

	for _, to := range msg.To {
		recipient, err := s.findRecipient(to)
		if err != nil {
			continue
		}

		err = email.Send(msg.From, recipient, msg.Subject, msg.Body)
		if err != nil {
			continue
		}
	}

	return err
}

func (s *Service) findRecipient(email string) (string, error) {
	name := email[:strings.Index(email, "@")]
	alias, err := s.GetAliasByName(name)
	if err != nil {
		return "", err
	}

	recipient, err := s.GetRecipient(context.Background(), alias.RecipientID)
	if err != nil {
		return "", err
	}

	return recipient.Email, nil
}

func parse(origin net.Addr, from string, to []string, data []byte) (Message, error) {
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
