package main

import (
	"log"

	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/cron"
	"ivpn.net/email/api/internal/repository"
	"ivpn.net/email/api/internal/service"
	"ivpn.net/email/api/internal/transport/api"
)

func Run() error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	// utils.NewLogger(cfg.API)

	db, err := repository.NewDB(cfg.DB)
	if err != nil {
		return err
	}

	redis, err := repository.NewRedis(cfg.Redis)
	if err != nil {
		return err
	}

	cron.New(db.Client)

	service := service.New(cfg, db, redis)

	err = api.Start(cfg.API, service, redis)
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
