package smpt

import (
	"log"

	"github.com/mhale/smtpd"
	"ivpn.net/email-service/config"
)

func Start(cfg config.Config) error {
	log.Printf("SMTP server starting on port %s", cfg.SMTPPort)
	srv := &smtpd.Server{
		Addr:         cfg.SMTPHost + ":" + cfg.SMTPPort,
		Handler:      InboundHandler,
		Hostname:     cfg.SMTPHostname,
		AuthRequired: true,
	}

	err := srv.ListenAndServe()
	return err
}
