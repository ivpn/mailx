package smpt

import (
	"log"
	"net"

	"github.com/mhale/smtpd"
	"ivpn.net/email/api/config"
)

type Service interface {
	ProcessMessage(net.Addr, string, []string, []byte) error
}

func Start(cfg config.SMTPConfig, service Service) error {
	log.Printf("SMTP server starting on :%s", cfg.Port)

	server := &smtpd.Server{
		Addr:         cfg.Host + ":" + cfg.Port,
		Handler:      service.ProcessMessage,
		Hostname:     cfg.Hostname,
		AuthRequired: false,
	}

	return server.ListenAndServe()
}
