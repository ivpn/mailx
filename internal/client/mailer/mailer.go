package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"net/mail"
	"strconv"

	"gopkg.in/gomail.v2"
	"ivpn.net/email-service/config"
)

var templateFS embed.FS

type Mailer struct {
	dialer *gomail.Dialer
	Sender string
}

type Msg struct {
	Subject string
	Body    string
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

func (mailer Mailer) Send(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mailer.Sender)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	m.AddAlternative("text/html", body)

	err := mailer.dialer.DialAndSend(m)
	if err != nil {
		return err
	}

	log.Println("Email sent successfully")
	return nil
}

func (mailer Mailer) Reply(from string, to string, data []byte) error {
	msg, err := parseMsg(data)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/plain", msg.Body)
	m.AddAlternative("text/html", msg.Body)

	err = mailer.dialer.DialAndSend(m)
	if err != nil {
		return err
	}

	log.Println("Email reply sent successfully")
	return nil
}

func (mailer Mailer) Forward(from string, to string, data []byte, templateFile string, templateData interface{}) error {
	msg, err := parseMsg(data)
	if err != nil {
		return err
	}

	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	header := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(header, "header", templateData)
	if err != nil {
		return err
	}

	headerHtml := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(headerHtml, "headerHtml", templateData)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/plain", header.String()+msg.Body)
	m.AddAlternative("text/html", headerHtml.String()+msg.Body)

	err = mailer.dialer.DialAndSend(m)
	if err != nil {
		return err
	}

	log.Println("Email forward sent successfully")
	return nil
}

func parseMsg(data []byte) (Msg, error) {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return Msg{}, err
	}

	subject := msg.Header.Get("Subject")
	buf := new(bytes.Buffer)
	buf.ReadFrom(msg.Body)
	body := buf.String()

	return Msg{
		Subject: subject,
		Body:    body,
	}, nil
}
