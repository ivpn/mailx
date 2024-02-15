package smpt

import (
	"log"
	"net"

	"github.com/mhale/smtpd"
	"ivpn.net/email-service/config"
)

type Service interface {
	ProcessMessage(net.Addr, string, []string, []byte) error
}

func Start(service Service) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	log.Printf("SMTP server starting on :%s", cfg.SMTP.Port)

	server := &smtpd.Server{
		Addr:         cfg.SMTP.Host + ":" + cfg.SMTP.Port,
		Handler:      service.ProcessMessage,
		Hostname:     cfg.SMTP.Hostname,
		AuthRequired: true,
	}

	return server.ListenAndServe()
}
