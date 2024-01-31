package server

import (
	"log"

	"github.com/mhale/smtpd"
	"ivpn.net/email-service/config"
	"ivpn.net/email-service/internal/mta/handler"
)

func Start(cfg config.Config) error {
	log.Printf("SMTP server starting on port %s", cfg.SMTPPort)
	srv := &smtpd.Server{
		Addr:         cfg.SMTPHost + ":" + cfg.SMTPPort,
		Handler:      handler.InboundHandler,
		Hostname:     cfg.SMTPHostname,
		AuthRequired: true,
	}

	err := srv.ListenAndServe()
	return err
}
