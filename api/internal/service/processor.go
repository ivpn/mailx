package service

import (
	"bytes"
	"context"
	"errors"
	"log"
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
	ErrNoRecipients         = errors.New("no verified recipients")
	ErrInactiveRecipient    = errors.New("inactive recipient")
)

type MsgType int

const (
	Reply MsgType = 2
	Send  MsgType = 3
)

type Msg struct {
	From     string
	FromName string
	To       []string
	Subject  string
	Body     string
	Type     MsgType
}

func (s *Service) ProcessMessage(data []byte) error {
	msg, err := parseMessage(data)
	if err != nil {
		log.Println("error parsing message", err)
		return err
	}

	for _, to := range msg.To {
		recipient, alias, relayType, err := s.findRecipient(msg.From, to, msg.Type)
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
			err = s.queueMessage(msg.From, msg.FromName, settings.FromName, recipient, data, alias, relayType)
			if err != nil {
				log.Println("error queueing message", err)
				return
			}

			s.saveMessage(alias, relayType, data)
		})
	}

	return err
}

func (s *Service) queueMessage(from string, fromName string, settingsFromName string, to string, data []byte, alias model.Alias, msgType model.MessageType) error {
	mailer := mailer.New(s.Cfg.SMTPClient)

	// Forward
	if msgType == model.Forward {
		templateData := map[string]interface{}{
			"alias": alias.Name,
			"from":  from,
		}
		generatedFrom := model.GenerateReplyTo(alias.Name, from)
		err := mailer.Forward(generatedFrom, fromName, to, data, "header.tmpl", templateData)
		if err != nil {
			log.Println("error forwarding message", err)
			return err
		}
	} else {
		// Reply | Send
		name := alias.FromName
		if name == "" {
			name = settingsFromName
		}

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

func (s *Service) findRecipient(from string, email string, msgType MsgType) (string, model.Alias, model.MessageType, error) {
	name, replyTo := model.ParseReplyTo(email)

	alias, err := s.GetAliasByName(name)
	if err != nil {
		return "", model.Alias{}, 0, err
	}

	if !alias.Enabled {
		// Block
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
		rcps, err := s.GetVerifiedRecipients(context.Background(), from, alias.UserID)
		if err != nil || len(rcps) == 0 {
			return "", model.Alias{}, 0, ErrNoRecipients
		}

		return replyTo, alias, model.MessageType(msgType), nil
	}

	r := strings.Split(alias.Recipients, ",")[0]

	return r, alias, model.Forward, nil
}

func parseMessage(data []byte) (Msg, error) {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return Msg{}, err
	}

	from := msg.Header.Get("From")
	to := strings.Split(msg.Header.Get("To"), ",")
	subject := msg.Header.Get("Subject")

	buf := new(bytes.Buffer)
	buf.ReadFrom(msg.Body)
	body := buf.String()
	msgType := Send

	fromHeader := msg.Header.Get("From")
	fromAddress, err := mail.ParseAddress(fromHeader)
	if err != nil {
		return Msg{}, err
	}
	fromName := fromAddress.Name

	if isReply(msg) {
		msgType = Reply
	}

	return Msg{
		From:     from,
		FromName: fromName,
		To:       to,
		Subject:  subject,
		Body:     body,
		Type:     msgType,
	}, nil
}

func isReply(m *mail.Message) bool {
	if m.Header.Get("In-Reply-To") != "" || m.Header.Get("References") != "" {
		return true
	}

	return false
}
