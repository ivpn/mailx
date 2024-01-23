package server

import (
	"log"

	"github.com/mhale/smtpd"
	"ivpn.net/email-service/services/inbound/config"
	"ivpn.net/email-service/services/inbound/handler"
)

func Start(cfg config.Config) error {
	log.Printf("Inbound server starting on %s", cfg.ServerPort)
	srv := &smtpd.Server{
		Addr:         cfg.ServerHost + ":" + cfg.ServerPort,
		Handler:      handler.InboundHandler,
		Hostname:     cfg.ServerHostname,
		AuthRequired: true,
	}

	return srv.ListenAndServe()
}
