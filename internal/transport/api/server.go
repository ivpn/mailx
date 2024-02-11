package api

import (
	"github.com/gofiber/fiber"
	"ivpn.net/email-service/config"
)

type Service interface {
	RecipientService
}

type Handler struct {
	Service Service
	Server  *fiber.App
}

func Start(cfg config.APIConfig, service Service) error {
	app := fiber.New()

	h := &Handler{
		Service: service,
		Server:  app,
	}

	h.SetupRoutes()

	return h.Server.Listen(cfg.Port)
}
