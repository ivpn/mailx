package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/mnako/letters"
	"github.com/yeo/parsemail"
	"gopkg.in/gomail.v2"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
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

func (mailer Mailer) ReplyLegacy(from string, name string, rcp model.Recipient, data []byte) error {
	var reader = bytes.NewReader(data)
	email, err := parsemail.Parse(reader)
	if err != nil {
		return err
	}

	if email.TextBody == "" {
		extractedTextBody, err := model.ExtractTextBody(data)
		if err != nil {
			log.Println("Error extracting text body:", err)
		} else {
			email.TextBody = extractedTextBody
		}
	}
	if email.HTMLBody == "" {
		extractedHTMLBody, err := model.ExtractHTMLBody(data)
		if err != nil {
			log.Println("Error extracting HTML body:", err)
		} else {
			email.HTMLBody = extractedHTMLBody
		}
	}
	if email.HTMLBody == "" {
		email.HTMLBody = model.PlainTextToHTML(email.TextBody)
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, name)
	m.SetHeader("To", rcp.Email)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/plain", utils.RemoveHeader(email.TextBody))
	m.AddAlternative("text/html", utils.RemoveHtmlHeader(email.HTMLBody))

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

	log.Printf("Email reply sent successfully, %s", email.MessageID)

	return nil
}

func (mailer Mailer) Reply(from string, name string, rcp model.Recipient, data []byte) error {
	reader := bytes.NewReader(data)
	email, err := letters.ParseEmail(reader)
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
	m.SetBody("text/plain", utils.RemoveHeader(email.Text))
	m.AddAlternative("text/html", utils.RemoveHtmlHeader(email.HTML))

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

func (mailer Mailer) ForwardLegacy(from string, name string, rcp model.Recipient, data []byte, templateFile string, templateData any) error {
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

	if email.TextBody == "" {
		extractedTextBody, err := model.ExtractTextBody(data)
		if err != nil {
			log.Println("Error extracting text body:", err)
		} else {
			email.TextBody = extractedTextBody
		}
	}
	if email.HTMLBody == "" {
		extractedHTMLBody, err := model.ExtractHTMLBody(data)
		if err != nil {
			log.Println("Error extracting HTML body:", err)
		} else {
			email.HTMLBody = extractedHTMLBody
		}
	}
	if email.HTMLBody == "" {
		email.HTMLBody = model.PlainTextToHTML(email.TextBody)
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, name)
	m.SetHeader("To", rcp.Email)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/plain", header.String()+email.TextBody)

	// PGP/Inline encryption
	if rcp.PGPEnabled && rcp.PGPKey != "" && rcp.PGPInline {
		pgp := crypto.PGP()
		publicKey, _ := crypto.NewKeyFromArmored(rcp.PGPKey)
		encHandle, _ := pgp.Encryption().Recipient(publicKey).New()
		pgpMessage, _ := encHandle.Encrypt([]byte(email.TextBody))
		armored, _ := pgpMessage.ArmorBytes()
		email.TextBody = string(armored)
		m.SetHeader("Content-Type", "text/plain")
		m.SetBody("text/plain", email.TextBody)
	} else {
		m.AddAlternative("text/html", headerHtml.String()+email.HTMLBody)
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

	// PGP/MIME encryption
	if rcp.PGPEnabled && rcp.PGPKey != "" && !rcp.PGPInline {
		var buf bytes.Buffer
		_, err = m.WriteTo(&buf)
		if err != nil {
			return err
		}

		pgp := crypto.PGP()
		publicKey, _ := crypto.NewKeyFromArmored(rcp.PGPKey)
		encHandle, _ := pgp.Encryption().Recipient(publicKey).New()
		pgpMessage, _ := encHandle.Encrypt(buf.Bytes())
		armored, _ := pgpMessage.ArmorBytes()

		msg := gomail.NewMessage()
		msg.SetAddressHeader("From", from, name)
		msg.SetHeader("To", rcp.Email)
		msg.SetHeader("Subject", email.Subject)
		msg.SetHeader("Content-Type", "multipart/encrypted; protocol=\"application/pgp-encrypted\"")
		msg.SetHeader("Content-Description", "OpenPGP encrypted message")
		msg.SetHeader("Content-Disposition", "inline; filename=\"encrypted.asc\"")
		msg.SetBody("application/pgp-encrypted", "Version: 1")
		msg.AddAlternative("application/octet-stream; name=\"encrypted.asc\"\r\n", string(armored))

		err = mailer.dialer.DialAndSend(msg)
		if err != nil {
			return err
		}

		log.Printf("PGP/MIME email forward sent successfully, %s", email.MessageID)
		return nil
	}

	err = mailer.dialer.DialAndSend(m)
	if err != nil {
		return err
	}

	log.Printf("Email forward sent successfully, %s", email.MessageID)
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
		pgp := crypto.PGP()
		publicKey, _ := crypto.NewKeyFromArmored(rcp.PGPKey)
		encHandle, _ := pgp.Encryption().Recipient(publicKey).New()
		pgpMessage, _ := encHandle.Encrypt([]byte(email.Text))
		armored, _ := pgpMessage.ArmorBytes()
		email.Text = string(armored)
		m.SetHeader("Content-Type", "text/plain")
		m.SetBody("text/plain", email.Text)
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
		var buf bytes.Buffer
		_, err = m.WriteTo(&buf)
		if err != nil {
			return err
		}

		pgp := crypto.PGP()
		publicKey, _ := crypto.NewKeyFromArmored(rcp.PGPKey)
		encHandle, _ := pgp.Encryption().Recipient(publicKey).New()
		pgpMessage, _ := encHandle.Encrypt(buf.Bytes())
		armored, _ := pgpMessage.ArmorBytes()

		msg := gomail.NewMessage()
		msg.SetAddressHeader("From", from, name)
		msg.SetHeader("To", rcp.Email)
		msg.SetHeader("Subject", email.Headers.Subject)
		msg.SetHeader("Content-Type", "multipart/encrypted; protocol=\"application/pgp-encrypted\"")
		msg.SetHeader("Content-Description", "OpenPGP encrypted message")
		msg.SetHeader("Content-Disposition", "inline; filename=\"encrypted.asc\"")
		msg.SetBody("application/pgp-encrypted", "Version: 1")
		msg.AddAlternative("application/octet-stream; name=\"encrypted.asc\"\r\n", string(armored))

		err = mailer.dialer.DialAndSend(msg)
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
