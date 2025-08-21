package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/mnako/letters"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
	"ivpn.net/email/api/internal/utils/gomail.v2"
)

//go:embed templates/*
var templateFS embed.FS

type Mailer struct {
	dialer     *gomail.Dialer
	Sender     string
	SenderName string
}

// #nosec G104
func New(cfg config.SMTPClientConfig) Mailer {
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		log.Println("Invalid SMTP port:", cfg.Port)
		return Mailer{
			dialer:     nil,
			Sender:     cfg.Sender,
			SenderName: cfg.SenderName,
		}
	}

	hosts := strings.Split(cfg.Host, ",")

	var dialer *gomail.Dialer
	for _, host := range hosts {
		host = strings.TrimSpace(host)
		if cfg.User == "" || cfg.Password == "" {
			dialer = &gomail.Dialer{Host: host, Port: port}
		} else {
			dialer = gomail.NewDialer(host, port, cfg.User, cfg.Password)
		}

		conn, err := dialer.Dial()
		if err != nil {
			log.Printf("Failed to connect to SMTP host: %s, trying next host if available. Error: %v\n", host, err)
			dialer = nil
		} else {
			conn.Close()
			break
		}
	}

	if dialer == nil {
		return Mailer{
			dialer:     nil,
			Sender:     cfg.Sender,
			SenderName: cfg.SenderName,
		}
	}

	return Mailer{
		dialer:     dialer,
		Sender:     cfg.Sender,
		SenderName: cfg.SenderName,
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

func (mailer Mailer) Reply(from string, name string, rcp model.Recipient, data []byte) error {
	reader := bytes.NewReader(data)
	email, err := letters.ParseEmail(reader)
	if err != nil {
		return err
	}

	email.Text = utils.RemoveHeader(email.Text)
	email.HTML = utils.RemoveHtmlHeader(email.HTML)

	if email.HTML == "" {
		email.HTML = model.PlainTextToHTML(email.Text)
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, name)
	m.SetHeader("To", rcp.Email)
	m.SetHeader("Subject", email.Headers.Subject)
	m.SetBody("text/plain", email.Text)
	m.AddAlternative("text/html", email.HTML)

	for _, a := range email.AttachedFiles {
		m.Attach(a.ContentDisposition.Params["filename"], gomail.SetCopyFunc(func(w io.Writer) error {
			_, err = w.Write(a.Data)
			return err
		}))
	}

	for _, f := range email.InlineFiles {
		m.Embed(
			f.ContentID,
			gomail.SetHeader(map[string][]string{
				"Content-ID":          {f.ContentID},
				"Content-Type":        {f.ContentType.ContentType},
				"Content-Disposition": {f.ContentDisposition.Params["type"] + "; filename=\"" + f.ContentDisposition.Params["filename"] + "\""},
			}),
			gomail.SetCopyFunc(func(w io.Writer) error {
				_, err = w.Write(f.Data)
				return err
			}),
		)
	}

	err = mailer.dialer.DialAndSend(m)
	if err != nil {
		return err
	}

	log.Printf("Email reply sent successfully, %s", email.Headers.MessageID)

	return nil
}

func (mailer Mailer) Forward(from string, name string, rcp model.Recipient, data []byte, templateFile string, templateData any) error {
	reader := bytes.NewReader(data)
	email, err := letters.ParseEmail(reader)
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

	if email.HTML == "" {
		email.HTML = model.PlainTextToHTML(email.Text)
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, name)
	m.SetHeader("To", rcp.Email)
	m.SetHeader("Subject", email.Headers.Subject)
	m.SetBody("text/plain", header.String()+email.Text)

	// PGP/Inline encryption
	if rcp.PGPEnabled && rcp.PGPKey != "" && rcp.PGPInline {
		armored, err := utils.EncryptWithPGPInline(email.Text, rcp.PGPKey)
		if err != nil {
			return err
		}
		m.SetHeader("Content-Type", "text/plain")
		m.SetBody("text/plain", armored)
	} else {
		m.AddAlternative("text/html", headerHtml.String()+email.HTML)
	}

	// PGPSignatures
	pgpSignatures, err := model.ExtractPGPSignatures(data)
	if err != nil {
		log.Println("Error extracting PGP signatures:", err)
	} else {
		for _, a := range pgpSignatures {
			m.Attach(a.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
				data, err := io.ReadAll(a.Data)
				if err != nil {
					return err
				}
				_, err = w.Write(data)
				return err
			}))
		}
	}

	// PGPKeys
	pgpKeys, err := model.ExtractPGPKeys(data)
	if err != nil {
		log.Println("Error extracting PGP keys:", err)
	} else {
		for _, a := range pgpKeys {
			m.Attach(a.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
				data, err := io.ReadAll(a.Data)
				if err != nil {
					return err
				}
				_, err = w.Write(data)
				return err
			}))
		}
	}

	for _, a := range email.AttachedFiles {
		m.Attach(a.ContentDisposition.Params["filename"], gomail.SetCopyFunc(func(w io.Writer) error {
			_, err = w.Write(a.Data)
			return err
		}))
	}

	for _, f := range email.InlineFiles {
		m.Embed(
			f.ContentID,
			gomail.SetHeader(map[string][]string{
				"Content-ID":          {f.ContentID},
				"Content-Type":        {f.ContentType.ContentType},
				"Content-Disposition": {f.ContentDisposition.Params["type"] + "; filename=\"" + f.ContentDisposition.Params["filename"] + "\""},
			}),
			gomail.SetCopyFunc(func(w io.Writer) error {
				_, err = w.Write(f.Data)
				return err
			}),
		)
	}

	// PGP/MIME encryption
	if rcp.PGPEnabled && rcp.PGPKey != "" && !rcp.PGPInline {
		em, err := utils.EncryptWithPGPMIME(m, from, name, email.Headers.Subject, rcp.Email, rcp.PGPKey)
		if err != nil {
			return err
		}

		err = mailer.dialer.DialAndSend(em)
		if err != nil {
			return err
		}

		log.Printf("PGP/MIME email forward sent successfully, %s", email.Headers.MessageID)
		return nil
	}

	err = mailer.dialer.DialAndSend(m)
	if err != nil {
		return err
	}

	log.Printf("Email forward sent successfully, %s", email.Headers.MessageID)
	return nil
}

func (mailer Mailer) SendTemplate(to string, subject string, templateFile string, templateData any) error {
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
