package api

import (
	"context"

	"github.com/gofiber/fiber"
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
	PostAlias(context.Context, model.Alias) error
	UpdateAlias(context.Context, string, model.Alias) error
	DeleteAlias(context.Context, string) error
}

func (h *Handler) GetAlias(c *fiber.Ctx) {
	id := c.Params("id")
	alias, err := h.Service.GetAlias(c.Context(), id)
	if err != nil {
		c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
		return
	}

	c.JSON(alias)
}

func (h *Handler) GetAliases(c *fiber.Ctx) {
	recipientID := c.Params("recipient_id")
	aliases, err := h.Service.GetAliases(c.Context(), recipientID)
	if err != nil {
		c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
		return
	}

	c.JSON(aliases)
}

func (h *Handler) PostAlias(c *fiber.Ctx) {
	var alias model.Alias
	err := c.BodyParser(&alias)
	if err != nil {
		c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
		return
	}

	err = h.Service.PostAlias(c.Context(), alias)
	if err != nil {
		c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
		return
	}

	c.Status(201).JSON(fiber.Map{
		"message": PostAliasSuccess,
	})
}

func (h *Handler) UpdateAlias(c *fiber.Ctx) {
	id := c.Params("id")
	var alias model.Alias
	err := c.BodyParser(&alias)
	if err != nil {
		c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
		return
	}

	err = h.Service.UpdateAlias(c.Context(), id, alias)
	if err != nil {
		c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
		return
	}

	c.Status(200).JSON(fiber.Map{
		"message": UpdateAliasSuccess,
	})
}

func (h *Handler) DeleteAlias(c *fiber.Ctx) {
	id := c.Params("id")
	err := h.Service.DeleteAlias(c.Context(), id)
	if err != nil {
		c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
		return
	}

	c.Status(200).JSON(fiber.Map{
		"message": DeleteAliasSuccess,
	})
}
