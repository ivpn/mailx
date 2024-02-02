package main

import (
	"log"

	"ivpn.net/email-service/config"
	"ivpn.net/email-service/internal/transport/smpt"
)

func Run() error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	err = smpt.Start(cfg)
	return err
}

func main() {
	err := Run()
	if err != nil {
		log.Println(err)
	}
}
