package service

import (
	"bytes"
	"context"
	"log"
	"net"
	"net/mail"
	"strings"
)

func (s *Service) ProcessMessage(origin net.Addr, from string, to []string, data []byte) error {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		log.Println("[Service] Drop message: " + err.Error())
	}

	log.Printf("[Service] From: %s", from)
	log.Printf("[Service] To: %s", to)
	log.Printf("[Service] Subject: %s", msg.Header.Get("Subject"))

	name := getName(to[0])
	alias, err := s.GetAliasByName(name)
	if err != nil {
		log.Println("[Service] Drop message: " + err.Error())
	}

	recipient, err := s.GetRecipient(context.Background(), alias.RecipientID)
	if err != nil {
		log.Println("[Service] Drop message: " + err.Error())
	}

	log.Println("[Service] Recipient: " + recipient.Email)

	return err
}

func getName(email string) string {
	return email[:strings.Index(email, "@")]
}
