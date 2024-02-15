package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email-service/config"
)

type Service interface {
	RecipientService
	AliasService
}

type Handler struct {
	Service Service
	Server  *fiber.App
}

func Start(service Service) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	log.Printf("API server starting on :%s", cfg.API.Port)

	app := fiber.New()

	h := &Handler{
		Service: service,
		Server:  app,
	}

	h.SetupRoutes()

	return h.Server.Listen(":" + cfg.API.Port)
}
