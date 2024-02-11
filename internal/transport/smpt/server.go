package smpt

import (
	"log"

	"github.com/mhale/smtpd"
	"ivpn.net/email-service/config"
)

func Start(cfg config.SMTPConfig) error {
	log.Printf("SMTP server starting on :%s", cfg.Port)

	server := &smtpd.Server{
		Addr:         cfg.Host + ":" + cfg.Port,
		Handler:      InboundHandler,
		Hostname:     cfg.Hostname,
		AuthRequired: true,
	}

	return server.ListenAndServe()
}
