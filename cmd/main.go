package main

import (
	"log"

	"ivpn.net/email-service/internal/repository"
	"ivpn.net/email-service/internal/service"
	"ivpn.net/email-service/internal/transport/api"
	"ivpn.net/email-service/internal/transport/smpt"
)

func main() {
	db, err := repository.New()
	if err != nil {
		log.Println(err)
	}

	svc := service.New(db)

	go func() {
		err := smpt.Start(svc)
		if err != nil {
			log.Println(err)
		}
	}()

	err = api.Start(svc)
	if err != nil {
		log.Println(err)
	}
}
