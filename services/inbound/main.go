package main

import (
	"log"

	"ivpn.net/email-service/services/inbound/config"
	"ivpn.net/email-service/services/inbound/server"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Println(err)
	}

	err = server.Start(cfg)
	if err != nil {
		log.Println(err)
	}
}
