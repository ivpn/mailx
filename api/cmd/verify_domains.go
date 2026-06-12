package main

import (
	"fmt"
	"log"
	"os"

	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/cron/jobs"
	"ivpn.net/email/api/internal/repository"
)

func runVerifyDomains() error {
	cfg := config.Config{
		API: config.APIConfig{
			TokenSecret: os.Getenv("TOKEN_SECRET"),
			Domains:     os.Getenv("DOMAINS"),
		},
		SMTPClient: config.SMTPClientConfig{
			DkimSelector: os.Getenv("SMTP_CLIENT_DKIM_SELECTOR"),
		},
		DB: dbConfigFromEnv(),
	}

	db, err := repository.NewDB(cfg.DB)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}

	log.Println("starting domain verification job")
	jobs.VerifyDomainsJob(cfg, db.Client)
	log.Println("domain verification job complete")
	return nil
}
