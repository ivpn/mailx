package mailer

import (
	"log"
	"strconv"

	"gopkg.in/gomail.v2"
	"ivpn.net/email-service/config"
)

func Send(from string, to string, subject string, body string) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(cfg.SMTPClient.Port)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.Dialer{
		Host:     cfg.SMTPClient.Host,
		Port:     port,
		Username: cfg.SMTPClient.User,
		Password: cfg.SMTPClient.Password,
	}
	err = d.DialAndSend(m)
	if err != nil {
		return err
	}

	log.Println("Email sent successfully")

	return nil
}
