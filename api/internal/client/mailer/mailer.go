package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"log"
	"strconv"

	"github.com/DusanKasan/parsemail"
	"gopkg.in/gomail.v2"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/model"
)

//go:embed templates/*
var templateFS embed.FS

type Mailer struct {
	dialer     *gomail.Dialer
	Sender     string
	SenderName string
}

func New(cfg config.SMTPClientConfig) Mailer {
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		log.Println(err)
	}

	if cfg.User == "" || cfg.Password == "" {
		return Mailer{
			dialer: &gomail.Dialer{Host: cfg.Host, Port: port},
			Sender: cfg.Sender,
		}
	}

	return Mailer{
		dialer: gomail.NewDialer(cfg.Host, port, cfg.User, cfg.Password),
		Sender: cfg.Sender,
	}
}

func (mailer Mailer) Send(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", mailer.Sender, mailer.SenderName)
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

func (mailer Mailer) Reply(from string, name string, to string, data []byte) error {
	var reader = bytes.NewReader(data)
	email, err := parsemail.Parse(reader)
	if err != nil {
		return err
	}

	if email.HTMLBody == "" {
		email.HTMLBody = model.PlainTextToHTML(email.TextBody)
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, name)
	m.SetHeader("To", to)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/plain", email.TextBody)
	m.AddAlternative("text/html", email.HTMLBody)

	for _, a := range email.Attachments {
		m.Attach(a.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
			data, err := io.ReadAll(a.Data)
			if err != nil {
				return err
			}
			_, err = w.Write(data)
			return err
		}))
	}

	for _, f := range email.EmbeddedFiles {
		m.Embed(f.CID, gomail.SetCopyFunc(func(w io.Writer) error {
			data, err := io.ReadAll(f.Data)
			if err != nil {
				return err
			}
			_, err = w.Write(data)
			return err
		}))
	}

	err = mailer.dialer.DialAndSend(m)
	if err != nil {
		return err
	}

	log.Println("Email reply sent successfully")
	return nil
}

func (mailer Mailer) Forward(from string, name string, to string, data []byte, templateFile string, templateData interface{}) error {
	var reader = bytes.NewReader(data)
	email, err := parsemail.Parse(reader)
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

	if email.HTMLBody == "" {
		email.HTMLBody = model.PlainTextToHTML(email.TextBody)
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, name)
	m.SetHeader("To", to)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/plain", header.String()+email.TextBody)
	m.AddAlternative("text/html", headerHtml.String()+email.HTMLBody)

	for _, a := range email.Attachments {
		m.Attach(a.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
			data, err := io.ReadAll(a.Data)
			if err != nil {
				return err
			}
			_, err = w.Write(data)
			return err
		}))
	}

	for _, f := range email.EmbeddedFiles {
		m.Embed(f.CID, gomail.SetCopyFunc(func(w io.Writer) error {
			data, err := io.ReadAll(f.Data)
			if err != nil {
				return err
			}
			_, err = w.Write(data)
			return err
		}))
	}

	err = mailer.dialer.DialAndSend(m)
	if err != nil {
		return err
	}

	log.Println("Email forward sent successfully")
	return nil
}

func (mailer Mailer) SendTemplate(to string, subject string, templateFile string, templateData interface{}) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", templateData)
	if err != nil {
		return err
	}

	bodyHtml := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(bodyHtml, "bodyHtml", templateData)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", mailer.Sender, mailer.SenderName)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body.String())
	m.AddAlternative("text/html", bodyHtml.String())

	err = mailer.dialer.DialAndSend(m)
	if err != nil {
		return err
	}

	log.Println("Email template sent successfully")
	return nil
}
