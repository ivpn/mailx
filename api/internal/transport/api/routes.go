package api

import (
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/swagger"
	"ivpn.net/email/api/config"
	_ "ivpn.net/email/api/docs"
	"ivpn.net/email/api/internal/middleware/auth"
)

func (h *Handler) SetupRoutes(cfg config.APIConfig) {
	h.Server.Use(auth.NewAPICORS(cfg))
	h.Server.Use(healthcheck.New())

	p := h.Server.Group("/v1/p")
	p.Use(limiter.New())
	p.Post("/register", h.Register)
	p.Post("/login", h.Login)

	sub := h.Server.Group("/v1/subscription/update")
	sub.Use(auth.NewPSK(cfg))
	sub.Use(auth.NewPSKCORS(cfg))
	sub.Put("", h.UpdateSubscription)

	v1 := h.Server.Group("/v1")
	v1.Use(auth.New(cfg, h.Cache))

	v1.Post("/user/sendotp", h.SendUserOTP)
	v1.Post("/user/activate", h.Activate)
	v1.Post("/user/logout", h.Logout)
	v1.Post("/user/delete", h.DeleteUser)
	v1.Get("/user", h.GetUser)
	v1.Get("/user/stats", h.GetUserStats)

	v1.Get("/subscription", h.GetSubscription)

	v1.Get("/settings", h.GetSettings)
	v1.Put("/settings", h.UpdateSettings)

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

	docs := h.Server.Group("/docs")
	docs.Use(auth.NewBasicAuth(cfg))
	docs.Get("/*", swagger.HandlerDefault)
}
