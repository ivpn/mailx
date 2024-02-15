package main

import (
	"log"

	"ivpn.net/email-service/internal/repository"
	"ivpn.net/email-service/internal/service"
	"ivpn.net/email-service/internal/transport/api"
	"ivpn.net/email-service/internal/transport/smpt"
)

func Run() error {
	db, err := repository.New()
	if err != nil {
		return err
	}

	service := service.New(db)

	go func() {
		smpt.Start(service)
	}()

	err = api.Start(service)
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
