package api

import (
	"github.com/gofiber/fiber/v2"
)

type ProcessorService interface {
	ProcessMessage([]byte) error
}

// @Summary Email handler
// @Description Handle incoming email
// @Tags email
// @Accept json
// @Produce json
// @Param email body string true "Email body"
// @Success 200 {string} string "OK"
// @Router /email [post]
func (h *Handler) HandleEmail(c *fiber.Ctx) error {
	err := h.Service.ProcessMessage(c.Body())
	if err != nil {
		return err
	}
	return c.SendString("OK")
}
