package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email-service/internal/middleware/auth"
	"ivpn.net/email-service/internal/model"
)

var (
	PostAliasSuccess   = "Alias created"
	UpdateAliasSuccess = "Alias updated"
	DeleteAliasSuccess = "Alias deleted"
)

type AliasService interface {
	GetAlias(context.Context, string) (model.Alias, error)
	GetAliases(context.Context, string) ([]model.Alias, error)
	PostAlias(context.Context, model.Alias, string) error
	UpdateAlias(context.Context, model.Alias) error
	DeleteAlias(context.Context, string) error
}

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

func (h *Handler) PostAlias(c *fiber.Ctx) error {
	req := AliasReq{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	alias := model.Alias{
		UserID:      auth.GetUserID(c),
		Description: req.Description,
		Enabled:     req.Enabled,
		Recipients:  req.Recipients,
		FromName:    req.FromName,
	}

	err = h.Service.PostAlias(c.Context(), alias, req.Format)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": PostAliasSuccess,
	})
}

func (h *Handler) UpdateAlias(c *fiber.Ctx) error {
	req := AliasReq{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	alias := model.Alias{
		Description: req.Description,
		Enabled:     req.Enabled,
		Recipients:  req.Recipients,
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
