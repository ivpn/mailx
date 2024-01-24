package main

import (
	"log"

	"ivpn.net/email-service/config"
	"ivpn.net/email-service/internal/mta/server"
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
