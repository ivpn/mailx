package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email-service/internal/middleware/auth"
	"ivpn.net/email-service/internal/model"
)

var (
	PostRecipientSuccess     = "Recipient created"
	ActivateRecipientSuccess = "Recipient activated"
	UpdateRecipientSuccess   = "Recipient updated"
	DeleteRecipientSuccess   = "Recipient deleted"
)

type RecipientService interface {
	GetRecipient(context.Context, string) (model.Recipient, error)
	GetRecipients(context.Context, string) ([]model.Recipient, error)
	PostRecipient(context.Context, model.Recipient) error
	UpdateRecipient(context.Context, model.Recipient) error
	DeleteRecipient(context.Context, string) error
	ActivateRecipient(context.Context, string, string) error
}

type RecipientRequest struct {
	Email string `json:"email"`
}

func (h *Handler) GetRecipient(c *fiber.Ctx) error {
	id := c.Params("id")
	recipient, err := h.Service.GetRecipient(c.Context(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(recipient)
}

func (h *Handler) GetRecipients(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	recipients, err := h.Service.GetRecipients(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(recipients)
}

func (h *Handler) PostRecipient(c *fiber.Ctx) error {
	req := RecipientRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	recipient := model.Recipient{
		UserID:   auth.GetUserID(c),
		Email:    req.Email,
		IsActive: false,
	}

	err = h.Service.PostRecipient(c.Context(), recipient)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": PostRecipientSuccess,
	})
}

func (h *Handler) DeleteRecipient(c *fiber.Ctx) error {
	ID := c.Params("id")
	err := h.Service.DeleteRecipient(c.Context(), ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": DeleteRecipientSuccess,
	})
}

func (h *Handler) ActivateRecipient(c *fiber.Ctx) error {
	ID := c.Params("id")

	req := ActivateRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = h.Service.ActivateRecipient(c.Context(), ID, req.OTP)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": ActivateRecipientSuccess,
	})
}
