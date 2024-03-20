package api

import (
	"ivpn.net/email-service/config"
	"ivpn.net/email-service/internal/middleware/auth"
)

func (h *Handler) SetupRoutes(cfg config.APIConfig) {
	h.Server.Post("/v1/register", h.Register)
	h.Server.Post("/v1/login", h.Login)

	sub := h.Server.Group("/v1/subscription/update")
	sub.Use(auth.NewPSK(cfg))
	sub.Use(auth.NewPSKCORS(cfg))
	sub.Put("", h.UpdateSubscription)

	v1 := h.Server.Group("/v1")
	v1.Use(auth.New(cfg))

	v1.Post("/user/sendotp", h.SendUserOTP)
	v1.Post("/user/activate", h.Activate)
	v1.Post("/user/logout", h.Logout)
	v1.Delete("/user/delete", h.DeleteUser)

	v1.Get("/subscription", h.GetSubscription)

	v1.Get("/recipient/:id", h.GetRecipient)
	v1.Get("/recipients", h.GetRecipients)
	v1.Post("/recipient", h.PostRecipient)
	v1.Post("/recipient/sendotp/:id", h.SendRecipientOTP)
	v1.Post("/recipient/activate/:id", h.ActivateRecipient)
	v1.Delete("/recipient/:id", h.DeleteRecipient)

	v1.Get("/alias/:id", h.GetAlias)
	v1.Get("/aliases", h.GetAliases)
	v1.Post("/alias", h.PostAlias)
	v1.Put("/alias/:id", h.UpdateAlias)
	v1.Delete("/alias/:id", h.DeleteAlias)
}
