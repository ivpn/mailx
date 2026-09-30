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
	UpdateRemoveHeaderSuccess  = "Mailx header setting updated successfully."
	ErrRemoveHeaderRequired    = "The remove_header field is required."
	ErrRemoveHeaderInvalid     = "The remove_header field must be true or false."
	UpdateLogIssuesSuccess     = "Diagnostics setting updated successfully."
	ErrLogIssuesRequired       = "The log_issues field is required."
	ErrLogIssuesInvalid        = "The log_issues field must be true or false."
)

type SettingsService interface {
	GetSettings(context.Context, string) (model.Settings, error)
	UpdateSettings(context.Context, model.Settings) error
	UpdateBlockTrackers(context.Context, string, bool) error
	UpdateRemoveHeader(context.Context, string, bool) error
	UpdateLogIssues(context.Context, string, bool) error
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
	req := BlockTrackersReq{}
	return h.updateBoolSetting(c, &req, func() *bool { return req.BlockTrackers }, boolSetting{
		errRequired: ErrBlockTrackersRequired,
		errInvalid:  ErrBlockTrackersInvalid,
		success:     UpdateBlockTrackersSuccess,
		update:      h.Service.UpdateBlockTrackers,
	})
}

// @Summary Update Mailx header
// @Description Turn the Mailx header in forwarded emails off (remove_header true) or on
// @Tags settings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param settings body RemoveHeaderReq true "Remove Mailx header"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Failure 422 {object} ErrorRes
// @Router /settings/remove-header [put]
func (h *Handler) UpdateRemoveHeader(c *fiber.Ctx) error {
	req := RemoveHeaderReq{}
	return h.updateBoolSetting(c, &req, func() *bool { return req.RemoveHeader }, boolSetting{
		errRequired: ErrRemoveHeaderRequired,
		errInvalid:  ErrRemoveHeaderInvalid,
		success:     UpdateRemoveHeaderSuccess,
		update:      h.Service.UpdateRemoveHeader,
	})
}

// @Summary Update diagnostics
// @Description Turn diagnostic logging of failed deliveries and forwarding issues on or off
// @Tags settings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param settings body LogIssuesReq true "Enable diagnostics"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Failure 422 {object} ErrorRes
// @Router /settings/log-issues [put]
func (h *Handler) UpdateLogIssues(c *fiber.Ctx) error {
	req := LogIssuesReq{}
	return h.updateBoolSetting(c, &req, func() *bool { return req.LogIssues }, boolSetting{
		errRequired: ErrLogIssuesRequired,
		errInvalid:  ErrLogIssuesInvalid,
		success:     UpdateLogIssuesSuccess,
		update:      h.Service.UpdateLogIssues,
	})
}

// boolSetting describes a single on/off setting updated through its own endpoint.
type boolSetting struct {
	errRequired string
	errInvalid  string
	success     string
	update      func(context.Context, string, bool) error
}

// updateBoolSetting parses req from the body and stores the value returned by value, which
// must be a required JSON boolean.
func (h *Handler) updateBoolSetting(c *fiber.Ctx, req any, value func() *bool, s boolSetting) error {
	userID := auth.GetUserID(c)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": s.errInvalid,
		})
	}
	v := value()
	if v == nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": s.errRequired,
		})
	}

	err := s.update(c.Context(), userID, *v)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": s.success,
	})
}
