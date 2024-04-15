package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/config"
)

type Service interface {
	RecipientService
	AliasService
	UserService
	SubscriptionService
	SettingsService
}

type Handler struct {
	Cfg     config.APIConfig
	Service Service
	Server  *fiber.App
}

func Start(cfg config.APIConfig, service Service) error {
	log.Printf("API server starting on :%s", cfg.Port)

	app := fiber.New()

	h := &Handler{
		Cfg:     cfg,
		Service: service,
		Server:  app,
	}

	h.SetupRoutes(cfg)

	return h.Server.Listen(":" + cfg.Port)
}
