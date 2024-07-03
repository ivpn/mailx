package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type ProcessorService interface {
	ProcessMessage([]byte) error
}

func (h *Handler) HandleEmail(c *fiber.Ctx) error {
	log.Println("Email received")

	h.Service.ProcessMessage(c.Body())

	return c.SendString("OK")
}
