package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
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

func Start(cfg config.APIConfig, service Service) error {
	log.Printf("API server starting on :%s", cfg.Port)

	app := fiber.New()
	app.Use(healthcheck.New())

	h := &Handler{
		Service: service,
		Server:  app,
	}

	h.SetupRoutes(cfg)

	return h.Server.Listen(":" + cfg.Port)
}
