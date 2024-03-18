package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email-service/internal/middleware/auth"
	"ivpn.net/email-service/internal/model"
)

type SubsctiptionService interface {
	GetSubscription(context.Context, string) (model.Subscription, error)
}

func (h *Handler) GetSubscription(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)

	sub, err := h.Service.GetSubscription(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(sub)
}
