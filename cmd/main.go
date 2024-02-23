package main

import (
	"log"

	"ivpn.net/email-service/internal/repository"
	"ivpn.net/email-service/internal/service"
	"ivpn.net/email-service/internal/transport/api"
	"ivpn.net/email-service/internal/transport/smpt"
)

func main() {
	db, err := repository.NewDB()
	if err != nil {
		log.Println(err)
	}

	service := service.New(db)

	go func() {
		err := smpt.Start(service)
		if err != nil {
			log.Println(err)
		}
	}()

	err = api.Start(service)
	if err != nil {
		log.Println(err)
	}
}
