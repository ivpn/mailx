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

// @Summary Get subscription
// @Description Get subscription
// @Tags subscription
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} model.Subscription
// @Failure 500 {object} ErrorRes
// @Router /v1/subscription [get]
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

// @Summary Update subscription
// @Description Update subscription
// @Tags subscription
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} SuccessRes
// @Failure 500 {object} ErrorRes
// @Router /v1/subscription/update [put]
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
