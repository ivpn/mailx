package service

import (
	"bytes"
	"log"
	"net"
	"net/mail"
)

func (s *Service) ProcessMessage(origin net.Addr, from string, to []string, data []byte) error {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		log.Print(err)
	}

	log.Printf("[Service] From: %s", from)
	log.Printf("[Service] To: %s", to)
	log.Printf("[Service] Subject: %s", msg.Header.Get("Subject"))

	return err
}
