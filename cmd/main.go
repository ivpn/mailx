package main

import (
	"log"

	"ivpn.net/email-service/config"
	"ivpn.net/email-service/internal/mta/server"
)

func Run() error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	err = server.Start(cfg)
	return err
}

func main() {
	err := Run()
	if err != nil {
		log.Println(err)
	}
}
