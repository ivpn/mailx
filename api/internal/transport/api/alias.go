package api

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
)

var (
	PostAliasSuccess   = "Alias created"
	UpdateAliasSuccess = "Alias updated"
	DeleteAliasSuccess = "Alias deleted"
	ErrInvalidDomain   = "Invalid domain"
	ErrUnverifiedRcp   = "Recipient not verified"
)

type AliasService interface {
	GetAlias(context.Context, string) (model.Alias, error)
	GetAliases(context.Context, string) ([]model.Alias, error)
	PostAlias(context.Context, model.Alias, string, string) error
	UpdateAlias(context.Context, model.Alias) error
	DeleteAlias(context.Context, string) error
}

// @Summary Get alias
// @Description Get alias by ID
// @Tags alias
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Alias ID"
// @Success 200 {object} model.Alias
// @Failure 500 {object} ErrorRes
// @Router /alias/{id} [get]
func (h *Handler) GetAlias(c *fiber.Ctx) error {
	id := c.Params("id")
	alias, err := h.Service.GetAlias(c.Context(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(alias)
}

// @Summary Get aliases
// @Description Get all aliases
// @Tags alias
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} model.Alias
// @Failure 500 {object} ErrorRes
// @Router /aliases [get]
func (h *Handler) GetAliases(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	aliases, err := h.Service.GetAliases(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(aliases)
}

// @Summary Create alias
// @Description Create alias
// @Tags alias
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body AliasReq true "Alias request"
// @Success 201 {object} SuccessRes
// @Failure 500 {object} ErrorRes
// @Router /alias [post]
func (h *Handler) PostAlias(c *fiber.Ctx) error {
	req := AliasReq{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	if !strings.Contains(h.Cfg.Domains, req.Domain) {
		return c.Status(500).JSON(fiber.Map{
			"error": ErrInvalidDomain,
		})
	}

	rcps, err := h.Service.GetVerifiedRecipients(c.Context(), req.Recipients)
	if err != nil || len(rcps) == 0 {
		return c.Status(500).JSON(fiber.Map{
			"error": ErrUnverifiedRcp,
		})
	}

	alias := model.Alias{
		UserID:      auth.GetUserID(c),
		Description: req.Description,
		Enabled:     req.Enabled,
		Recipients:  model.GetEmails(rcps),
		FromName:    req.FromName,
	}

	err = h.Service.PostAlias(c.Context(), alias, req.Format, req.Domain)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": PostAliasSuccess,
	})
}

// @Summary Update alias
// @Description Update alias
// @Tags alias
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Alias ID"
// @Param body body AliasReq true "Alias request"
// @Success 200 {object} SuccessRes
// @Failure 500 {object} ErrorRes
// @Router /alias/{id} [put]
func (h *Handler) UpdateAlias(c *fiber.Ctx) error {
	req := AliasReq{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	rcps, err := h.Service.GetVerifiedRecipients(c.Context(), req.Recipients)
	if err != nil || len(rcps) == 0 {
		return c.Status(500).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	alias := model.Alias{
		Description: req.Description,
		Enabled:     req.Enabled,
		Recipients:  model.GetEmails(rcps),
		FromName:    req.FromName,
	}
	alias.ID = c.Params("id")

	err = h.Service.UpdateAlias(c.Context(), alias)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": UpdateAliasSuccess,
	})
}

// @Summary Delete alias
// @Description Delete alias
// @Tags alias
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Alias ID"
// @Success 200 {object} SuccessRes
// @Failure 500 {object} ErrorRes
// @Router /alias/{id} [delete]
func (h *Handler) DeleteAlias(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.Service.DeleteAlias(c.Context(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": DeleteAliasSuccess,
	})
}
