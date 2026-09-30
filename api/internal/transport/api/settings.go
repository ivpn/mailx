package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
)

var (
	UpdateSettingsSuccess      = "Settings updated successfully."
	UpdateBlockTrackersSuccess = "Block email trackers setting updated successfully."
	ErrBlockTrackersRequired   = "The block_trackers field is required."
	ErrBlockTrackersInvalid    = "The block_trackers field must be true or false."
)

type SettingsService interface {
	GetSettings(context.Context, string) (model.Settings, error)
	UpdateSettings(context.Context, model.Settings) error
	UpdateBlockTrackers(context.Context, string, bool) error
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
		UserID:       userID,
		Domain:       req.Domain,
		Recipient:    req.Recipient,
		FromName:     req.FromName,
		AliasFormat:  req.AliasFormat,
		LogIssues:    req.LogIssues,
		RemoveHeader: req.RemoveHeader,
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

// @Summary Update block email trackers
// @Description Turn removal of tracking pixels from forwarded emails on or off
// @Tags settings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param settings body BlockTrackersReq true "Block email trackers"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Failure 422 {object} ErrorRes
// @Router /settings/block-trackers [put]
func (h *Handler) UpdateBlockTrackers(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)

	req := BlockTrackersReq{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": ErrBlockTrackersInvalid,
		})
	}
	if req.BlockTrackers == nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": ErrBlockTrackersRequired,
		})
	}

	err := h.Service.UpdateBlockTrackers(c.Context(), userID, *req.BlockTrackers)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": UpdateBlockTrackersSuccess,
	})
}
