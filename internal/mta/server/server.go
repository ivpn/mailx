package server

import (
	"log"

	"github.com/mhale/smtpd"
	"ivpn.net/email-service/config"
	"ivpn.net/email-service/internal/mta/handler"
)

func Start(cfg config.Config) error {
	log.Printf("Inbound server starting on %s", cfg.MTAPort)
	srv := &smtpd.Server{
		Addr:         cfg.MTAHost + ":" + cfg.MTAPort,
		Handler:      handler.InboundHandler,
		Hostname:     cfg.MTAHostname,
		AuthRequired: true,
	}

	return srv.ListenAndServe()
}
