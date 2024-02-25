package mailer

import (
	"log"
	"strconv"

	"gopkg.in/gomail.v2"
	"ivpn.net/email-service/config"
)

type Mailer struct {
	dialer *gomail.Dialer
	Sender string
}

func New(cfg config.SMTPClientConfig) Mailer {
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		log.Println(err)
	}

	return Mailer{
		dialer: gomail.NewDialer(cfg.Host, port, cfg.User, cfg.Password),
		Sender: cfg.Sender,
	}
}

func (m Mailer) Send(to string, subject string, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.Sender)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	err := m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	log.Println("Email sent successfully")
	return nil
}
