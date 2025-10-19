package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
)

var (
	UpdateSubscriptionSuccess = "Subscription updated successfully."
	AddSubscriptionSuccess    = "Subscription added successfully."
	InvalidPASessionId        = "Invalid pre-auth session ID."
)

type SubscriptionService interface {
	GetSubscription(context.Context, string) (model.Subscription, error)
	UpdateSubscription(context.Context, model.Subscription, string, string, string) error
	AddPASession(context.Context, model.PASession) error
	RotatePASessionId(context.Context, string) (string, error)
}

// @Summary Get subscription
// @Description Get subscription
// @Tags subscription
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} model.Subscription
// @Failure 400 {object} ErrorRes
// @Router /sub [get]
func (h *Handler) GetSubscription(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)

	sub, err := h.Service.GetSubscription(c.Context(), userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
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
// @Param body body SubscriptionReq true "Subscription request"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /subscription/update [put]
func (h *Handler) UpdateSubscription(c *fiber.Ctx) error {
	req := SubscriptionReq{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	err = h.Validator.Struct(req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	sub := model.Subscription{}
	sub.ID = req.ID

	err = h.Service.UpdateSubscription(c.Context(), sub, req.SubID, req.PreauthID, req.PreauthTokenHash)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": UpdateSubscriptionSuccess,
	})
}

// @Summary Add pre-auth session
// @Description Add pre-auth session
// @Tags subscription
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body PASessionReq true "Pre-auth session request"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /sub/session [post]
func (h *Handler) AddPASession(c *fiber.Ctx) error {
	req := PASessionReq{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	err = h.Validator.Struct(req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	paSession := model.PASession{
		ID:        req.ID,
		PreAuthID: req.PreAuthID,
		Token:     req.Token,
	}

	err = h.Service.AddPASession(c.Context(), paSession)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return nil
}

// @Summary Rotate pre-auth session ID
// @Description Rotate pre-auth session ID
// @Tags subscription
// @Accept json
// @Produce json
// @Param body body RotatePASessionReq true "Rotate pre-auth session request"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /rotatepasession [put]
func (h *Handler) RotatePASession(c *fiber.Ctx) error {
	req := RotatePASessionReq{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": InvalidPASessionId,
		})
	}

	err = h.Validator.Struct(req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": InvalidPASessionId,
		})
	}

	newID, err := h.Service.RotatePASessionId(c.Context(), req.ID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": InvalidPASessionId,
		})
	}

	c.Cookie(auth.NewCookiePASession(newID))

	return c.SendStatus(fiber.StatusOK)
}
