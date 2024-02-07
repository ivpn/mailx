package main

import (
	"log"
	"strconv"

	"gopkg.in/gomail.v2"
	"ivpn.net/email-service/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Panic(err)
	}

	send(cfg.SMTP)
}

func send(cfg config.SMTPConfig) {
	host := cfg.Host
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		log.Panic(err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "from@example.com")
	m.SetHeader("To", "to@example.com")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/plain", "Hello!")

	d := gomail.Dialer{Host: host, Port: port}
	if err := d.DialAndSend(m); err != nil {
		log.Panic(err)
	}

	log.Println("Sent email")
}
