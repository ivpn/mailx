package api

import (
	"context"
	"log"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/utils"
)

type Service interface {
	RecipientService
	AliasService
	UserService
	SubscriptionService
	SettingsService
	ProcessorService
	SessionService
	CredentialService
	BounceService
	DiscardService
}

type Handler struct {
	Cfg       config.APIConfig
	Service   Service
	Server    *fiber.App
	Validator utils.Validator
	Cache     Cache
	WebAuthn  *webauthn.WebAuthn
}

type Cache interface {
	Get(context.Context, string) (string, error)
}

func Start(cfg config.APIConfig, service Service, cache Cache) error {
	log.Printf("API server starting on :%s", cfg.Port)

	app := fiber.New()

	h := &Handler{
		Cfg:       cfg,
		Service:   service,
		Server:    app,
		Validator: utils.NewValidator(),
		Cache:     cache,
		WebAuthn:  auth.NewWebAuthn(cfg),
	}

	h.SetupRoutes(cfg)

	return h.Server.Listen(":" + cfg.Port)
}
