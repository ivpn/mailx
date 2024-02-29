package api

import (
	"ivpn.net/email-service/config"
	"ivpn.net/email-service/internal/middleware/auth"
)

func (h *Handler) SetupRoutes(cfg config.APIConfig) {
	v1 := h.Server.Group("/v1")

	v1.Post("/register", h.Register)
	v1.Post("/login", h.Login)

	auth := auth.New(cfg)

	user := v1.Group("/user")
	user.Use(auth)
	user.Post("/activate", h.Activate)
	user.Post("/logout", h.Logout)

	recipient := v1.Group("/recipient")
	recipient.Use(auth)
	recipient.Get("/:id", h.GetRecipient)
	recipient.Get("/:user_id", h.GetRecipients)
	recipient.Post("/", h.PostRecipient)
	recipient.Put("/", h.UpdateRecipient)
	recipient.Delete("/:id", h.DeleteRecipient)
	recipient.Get("/verify/:id/:verification", h.VerifyRecipient)

	alias := v1.Group("/alias")
	alias.Use(auth)
	alias.Get("/:id", h.GetAlias)
	alias.Get("/:user_id", h.GetAliases)
	alias.Post("", h.PostAlias)
	alias.Put("", h.UpdateAlias)
	alias.Delete("/:id", h.DeleteAlias)
}
