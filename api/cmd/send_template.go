package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/client/mailer"
)

func smtpConfigFromEnv() config.SMTPClientConfig {
	return config.SMTPClientConfig{
		Host:        os.Getenv("SMTP_CLIENT_HOST"),
		Port:        os.Getenv("SMTP_CLIENT_PORT"),
		User:        os.Getenv("SMTP_CLIENT_USER"),
		Password:    os.Getenv("SMTP_CLIENT_PASSWORD"),
		Sender:      os.Getenv("SMTP_CLIENT_SENDER"),
		SenderName:  os.Getenv("SMTP_CLIENT_SENDER_NAME"),
		Report:      os.Getenv("SMTP_CLIENT_REPORT"),
		TokenSecret: os.Getenv("TOKEN_SECRET"),
	}
}

func runSendTemplate(args []string) error {
	fs := flag.NewFlagSet("send-template", flag.ContinueOnError)
	tmpl := fs.String("template", "", "Template filename (e.g. expiring_beta.tmpl)")
	to := fs.String("to", "", "Comma-separated list of recipient email addresses")
	subject := fs.String("subject", "", "Email subject")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *tmpl == "" {
		fs.Usage()
		return fmt.Errorf("--template is required")
	}
	if *to == "" {
		fs.Usage()
		return fmt.Errorf("--to is required")
	}
	if *subject == "" {
		fs.Usage()
		return fmt.Errorf("--subject is required")
	}

	recipients := strings.Split(*to, ",")
	for i, r := range recipients {
		recipients[i] = strings.TrimSpace(r)
	}

	cfg := smtpConfigFromEnv()
	m := mailer.New(cfg)

	var failed int
	for _, email := range recipients {
		if email == "" {
			continue
		}
		err := m.SendTemplate(email, *subject, *tmpl, nil)
		if err != nil {
			log.Printf("failed to send to %s: %s", email, err.Error())
			failed++
		} else {
			log.Printf("sent to %s", email)
		}
	}

	if failed > 0 {
		return fmt.Errorf("%d recipient(s) failed", failed)
	}
	return nil
}
