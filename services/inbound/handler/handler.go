package handler

import (
	"bytes"
	"log"
	"net"
	"net/mail"
)

func InboundHandler(origin net.Addr, from string, to []string, data []byte) error {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		log.Print(err)
	}

	log.Printf("From: %s", from)
	log.Printf("To: %s", to)
	log.Printf("Subject: %s", msg.Header.Get("Subject"))

	return err
}
