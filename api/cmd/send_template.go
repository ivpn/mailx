package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/client/mailer"
	"ivpn.net/email/api/internal/repository"
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

func dbConfigFromEnv() config.DBConfig {
	return config.DBConfig{
		Hosts:    strings.Split(os.Getenv("DB_HOSTS"), ","),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}
}

func runSendTemplateManaged(args []string) error {
	fs := flag.NewFlagSet("send-template-managed", flag.ContinueOnError)
	tmpl := fs.String("template", "", "Template filename (e.g. expiring_beta.tmpl)")
	subject := fs.String("subject", "", "Email subject")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *tmpl == "" {
		fs.Usage()
		return fmt.Errorf("--template is required")
	}
	if *subject == "" {
		fs.Usage()
		return fmt.Errorf("--subject is required")
	}

	db, err := repository.NewDB(dbConfigFromEnv())
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}

	emails, err := db.GetManagedUserEmails(context.Background())
	if err != nil {
		return fmt.Errorf("querying managed users: %w", err)
	}

	if len(emails) == 0 {
		log.Println("no managed users found")
		return nil
	}

	log.Printf("found %d managed user(s)", len(emails))

	m := mailer.New(smtpConfigFromEnv())

	var failed int
	for _, email := range emails {
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
