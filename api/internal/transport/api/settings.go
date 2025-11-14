package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
)

var (
	UpdateSettingsSuccess = "Settings updated successfully."
)

type SettingsService interface {
	GetSettings(context.Context, string) (model.Settings, error)
	UpdateSettings(context.Context, model.Settings) error
}

// @Summary Get settings
// @Description Get settings
// @Tags settings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} model.Settings
// @Failure 400 {object} ErrorRes
// @Router /settings [get]
func (h *Handler) GetSettings(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)

	settings, err := h.Service.GetSettings(c.Context(), userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(settings)
}

// @Summary Update settings
// @Description Update settings
// @Tags settings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param settings body SettingsReq true "Settings"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /settings [put]
func (h *Handler) UpdateSettings(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)

	req := SettingsReq{}
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

	settings := model.Settings{
		UserID:      userID,
		Domain:      req.Domain,
		Recipient:   req.Recipient,
		FromName:    req.FromName,
		AliasFormat: req.AliasFormat,
		LogBounce:   req.LogBounce,
		LogDiscard:  req.LogDiscard,
	}
	settings.ID = req.ID

	err = h.Service.UpdateSettings(c.Context(), settings)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": UpdateSettingsSuccess,
	})
}
