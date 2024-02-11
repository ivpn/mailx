package main

import (
	"log"

	"ivpn.net/email-service/config"
	"ivpn.net/email-service/internal/repository"
	"ivpn.net/email-service/internal/service"
	"ivpn.net/email-service/internal/transport/api"
	"ivpn.net/email-service/internal/transport/smpt"
)

func Run() error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	go func() {
		smpt.Start(cfg.SMTP)
	}()

	db, err := repository.New(cfg.DB)
	if err != nil {
		return err
	}

	service := service.New(db)

	err = api.Start(cfg.API, service)
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
