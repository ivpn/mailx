package api

import (
	"context"

	"github.com/gofiber/fiber"
	"ivpn.net/email-service/internal/model"
)

type RecipientService interface {
	GetRecipient(context.Context, string) (error, model.Recipient)
	GetRecipients(context.Context, string) ([]model.Recipient, error)
	PostRecipient(context.Context, model.Recipient) error
	UpdateRecipient(context.Context, model.Recipient) error
	DeleteRecipient(context.Context, string) error
	VerifyRecipient(context.Context, string, string) (model.Recipient, error)
}

func (h *Handler) GetRecipient(c *fiber.Ctx) {
	id := c.Params("id")
	err, model := h.Service.GetRecipient(c.Context(), id)
	if err != nil {
		c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
		return
	}

	c.JSON(model)
}

func (h *Handler) GetRecipients(c *fiber.Ctx) {
	userID := c.Params("userID")
	recipients, err := h.Service.GetRecipients(c.Context(), userID)
	if err != nil {
		c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
		return
	}

	c.JSON(recipients)
}

func (h *Handler) PostRecipient(c *fiber.Ctx) {
	var recipient model.Recipient
	err := c.BodyParser(&recipient)
	if err != nil {
		c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
		return
	}

	err = h.Service.PostRecipient(c.Context(), recipient)
	if err != nil {
		c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
		return
	}

	c.Status(201).JSON(fiber.Map{
		"message": "Recipient created",
	})
}

func (h *Handler) UpdateRecipient(c *fiber.Ctx) {
	var recipient model.Recipient
	err := c.BodyParser(&recipient)
	if err != nil {
		c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
		return
	}

	err = h.Service.UpdateRecipient(c.Context(), recipient)
	if err != nil {
		c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
		return
	}

	c.Status(200).JSON(fiber.Map{
		"message": "Recipient updated",
	})
}

func (h *Handler) DeleteRecipient(c *fiber.Ctx) {
	id := c.Params("id")
	err := h.Service.DeleteRecipient(c.Context(), id)
	if err != nil {
		c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
		return
	}

	c.Status(200).JSON(fiber.Map{
		"message": "Recipient deleted",
	})
}

func (h *Handler) VerifyRecipient(c *fiber.Ctx) {
	id := c.Params("id")
	verification := c.Params("verification")
	model, err := h.Service.VerifyRecipient(c.Context(), id, verification)
	if err != nil {
		c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
		return
	}

	c.JSON(model)
}
