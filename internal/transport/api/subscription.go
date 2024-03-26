package api

import (
	"context"

	"github.com/araddon/dateparse"
	"github.com/gofiber/fiber/v2"
	"ivpn.net/email-service/internal/middleware/auth"
	"ivpn.net/email-service/internal/model"
)

var (
	UpdateSubscriptionSuccess = "Subscription updated"
)

type SubscriptionService interface {
	GetSubscription(context.Context, string) (model.Subscription, error)
	UpdateSubscription(context.Context, model.Subscription) error
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

func (h *Handler) UpdateSubscription(c *fiber.Ctx) error {
	req := SubscriptionReq{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	activeUntil, err := dateparse.ParseAny(req.ActiveUntil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	sub := model.Subscription{}
	sub.ID = req.ID
	sub.ActiveUntil = activeUntil

	err = h.Service.UpdateSubscription(c.Context(), sub)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": UpdateSubscriptionSuccess,
	})
}
