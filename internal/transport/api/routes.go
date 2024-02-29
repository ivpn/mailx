package api

import (
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"ivpn.net/email-service/config"
	"ivpn.net/email-service/internal/middleware/auth"
)

func (h *Handler) SetupRoutes(cfg config.APIConfig) {
	v1 := h.Server.Group("/v1")

	v1.Use(encryptcookie.New(encryptcookie.Config{
		Key: cfg.CookieSecret,
	}))
	v1.Use(auth.New(cfg))

	v1.Post("/register", h.Register)
	v1.Post("/activate", h.Activate)
	v1.Post("/login", h.Login)
	v1.Post("/logout", h.Logout)

	v1.Get("/recipient/:id", h.GetRecipient)
	v1.Get("/recipients/:user_id", h.GetRecipients)
	v1.Post("/recipient", h.PostRecipient)
	v1.Put("/recipient", h.UpdateRecipient)
	v1.Delete("/recipient/:id", h.DeleteRecipient)
	v1.Get("/recipient/verify/:id/:verification", h.VerifyRecipient)

	v1.Get("/alias/:id", h.GetAlias)
	v1.Get("/aliases/:user_id", h.GetAliases)
	v1.Post("/alias", h.PostAlias)
	v1.Put("/alias", h.UpdateAlias)
	v1.Delete("/alias/:id", h.DeleteAlias)
}
