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
		return err
	}

	log.Printf("[Service] From: %s", from)
	log.Printf("[Service] To: %s", to)
	log.Printf("[Service] Subject: %s", msg.Header.Get("Subject"))

	for _, to := range to {
		name := getName(to)
		alias, err := s.GetAliasByName(name)
		if err != nil {
			log.Println("[Service] Drop message: " + err.Error())
			continue
		}

		recipient, err := s.GetRecipient(context.Background(), alias.RecipientID)
		if err != nil {
			log.Println("[Service] Drop message: " + err.Error())
			continue
		}

		log.Println("[Service] Recipient: " + recipient.Email)
	}

	return err
}

func getName(email string) string {
	return email[:strings.Index(email, "@")]
}
