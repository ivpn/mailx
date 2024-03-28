package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"strconv"

	"gopkg.in/gomail.v2"
	"ivpn.net/email-service/config"
)

var templateFS embed.FS

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
	msg.AddAlternative("text/html", body)

	err := m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	log.Println("Email sent successfully")
	return nil
}

func (m Mailer) Forward(to string, subject string, body string, templateFile string, data interface{}) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	header := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(header, "header", data)
	if err != nil {
		return err
	}

	headerHtml := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(headerHtml, "headerHtml", data)
	if err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", m.Sender)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", header.String()+body)
	msg.AddAlternative("text/html", headerHtml.String()+body)

	err = m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	log.Println("Email (Forward) sent successfully")
	return nil
}
