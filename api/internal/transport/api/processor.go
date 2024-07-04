package api

import (
	"github.com/gofiber/fiber/v2"
)

type ProcessorService interface {
	ProcessMessage([]byte) error
}

func (h *Handler) HandleEmail(c *fiber.Ctx) error {
	h.Service.ProcessMessage(c.Body())
	return c.SendString("OK")
}
