package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email-service/internal/middleware/auth"
	"ivpn.net/email-service/internal/model"
)

var (
	UpdateSettingsSuccess = "Settings updated"
)

type SettingsService interface {
	GetSettings(context.Context, string) (model.Settings, error)
	UpdateSettings(context.Context, model.Settings) error
}

func (h *Handler) GetSettings(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)

	settings, err := h.Service.GetSettings(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(settings)
}

func (h *Handler) UpdateSettings(c *fiber.Ctx) error {
	req := SettingsReq{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	settings := model.Settings{
		Domain:    req.Domain,
		Recipient: req.Recipient,
		FromName:  req.FromName,
	}
	settings.ID = req.ID

	err = h.Service.UpdateSettings(c.Context(), settings)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": UpdateSettingsSuccess,
	})
}
