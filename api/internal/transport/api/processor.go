package api

import (
	"github.com/gofiber/fiber/v2"
)

type ProcessorService interface {
	ProcessMessage(data []byte, envelopeRecipient string) error
}

// @Summary Email handler
// @Description Handle incoming email
// @Tags email
// @Accept json
// @Produce json
// @Param email body string true "Email body"
// @Param recipient query string false "Envelope recipient for this specific delivery (Postfix ${recipient})"
// @Success 200 {string} string "OK"
// @Router /email [post]
func (h *Handler) HandleEmail(c *fiber.Ctx) error {
	err := h.Service.ProcessMessage(c.Body(), c.Query("recipient"))
	if err != nil {
		// TEMPORARY failure → Postfix should retry
		return c.Status(fiber.StatusServiceUnavailable).SendString("temporary failure")
	}

	return c.SendString("OK")
}
