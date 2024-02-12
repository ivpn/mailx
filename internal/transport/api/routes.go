package api

import "github.com/gofiber/fiber/v2"

func (h *Handler) SetupRoutes() {
	v1 := h.Server.Group("/v1")

	v1.Get("/health", HealthCheck)

	v1.Get("/recipient/:id", h.GetRecipient)
	v1.Get("/recipients/:user_id", h.GetRecipients)
	v1.Post("/recipient", h.PostRecipient)
	v1.Put("/recipient", h.UpdateRecipient)
	v1.Delete("/recipient/:id", h.DeleteRecipient)
	v1.Get("/recipient/verify/:id/:verification", h.VerifyRecipient)

	v1.Get("/alias/:id", h.GetAlias)
	v1.Get("/aliases/:recipient_id", h.GetAliases)
	v1.Post("/alias", h.PostAlias)
	v1.Put("/alias", h.UpdateAlias)
	v1.Delete("/alias/:id", h.DeleteAlias)
}

func HealthCheck(c *fiber.Ctx) error {
	return c.SendString("OK")
}
