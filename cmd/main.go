package main

import (
	"log"

	"ivpn.net/email-service/config"
	"ivpn.net/email-service/internal/repository"
	"ivpn.net/email-service/internal/transport/smpt"
)

func Run() error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	_, err = repository.New(cfg.DB)
	if err != nil {
		return err
	}

	err = smpt.Start(cfg.SMTP)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := Run()
	if err != nil {
		log.Println(err)
	}
}
