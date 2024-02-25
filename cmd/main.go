package main

import (
	"log"

	"ivpn.net/email-service/config"
	"ivpn.net/email-service/internal/repository"
	"ivpn.net/email-service/internal/service"
	"ivpn.net/email-service/internal/transport/api"
	"ivpn.net/email-service/internal/transport/smpt"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Println(err)
	}

	db, err := repository.NewDB()
	if err != nil {
		log.Println(err)
	}

	redis, err := repository.NewRedis()
	if err != nil {
		log.Println(err)
	}

	service := service.New(cfg, db, redis)

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
