package api

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
)

var (
	PostRecipientSuccess     = "Recipient added successfully."
	ActivateRecipientSuccess = "Recipient activated successfully."
	UpdateRecipientSuccess   = "Recipient updated successfully."
	DeleteRecipientSuccess   = "Recipient deleted successfully."
)

type RecipientService interface {
	GetRecipient(context.Context, string, string) (model.Recipient, error)
	GetRecipients(context.Context, string) ([]model.Recipient, error)
	GetVerifiedRecipients(context.Context, string, string) ([]model.Recipient, error)
	PostRecipient(context.Context, model.Recipient) error
	SendRecipientOTP(context.Context, string, string) error
	UpdateRecipient(context.Context, model.Recipient) error
	ActivateRecipient(context.Context, string, string, string) error
	DeleteRecipient(context.Context, string, string) error
}

// @Summary Get recipient
// @Description Get recipient by ID
// @Tags recipient
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Recipient ID"
// @Success 200 {object} model.Recipient
// @Failure 400 {object} ErrorRes
// @Router /recipient/{id} [get]
func (h *Handler) GetRecipient(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	id := c.Params("id")
	recipient, err := h.Service.GetRecipient(c.Context(), id, userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(recipient)
}

// @Summary Get recipients
// @Description Get all recipients
// @Tags recipient
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} model.Recipient
// @Failure 400 {object} ErrorRes
// @Router /recipients [get]
func (h *Handler) GetRecipients(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	rcps, err := h.Service.GetRecipients(c.Context(), userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	for i, rcp := range rcps {
		if rcp.PGPKey != "" {
			rcps[i].PGPKey = utils.HashPGPKey(rcp.PGPKey)
		}
	}

	return c.JSON(rcps)
}

// @Summary Create recipient
// @Description Create recipient
// @Tags recipient
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body EmailReq true "Recipient request"
// @Success 201 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /recipient [post]
func (h *Handler) PostRecipient(c *fiber.Ctx) error {
	req := EmailReq{}
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

	recipient := model.Recipient{
		UserID:   auth.GetUserID(c),
		Email:    req.Email,
		IsActive: false,
	}

	err = h.Service.PostRecipient(c.Context(), recipient)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": PostRecipientSuccess,
	})
}

// @Summary Update recipient
// @Description Update recipient
// @Tags recipient
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body RecipientReq true "Recipient request"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /recipient [put]
func (h *Handler) UpdateRecipient(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)

	req := RecipientReq{}
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

	rcp, err := h.Service.GetRecipient(c.Context(), req.ID, userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if req.PGPKey == "" || strings.HasPrefix(req.PGPKey, "-----BEGIN PGP PUBLIC KEY BLOCK-----") {
		rcp.PGPKey = req.PGPKey
	}

	rcp.PGPEnabled = req.PGPEnabled
	rcp.PGPInline = req.PGPInline

	err = h.Service.UpdateRecipient(c.Context(), rcp)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": UpdateRecipientSuccess,
	})
}

// @Summary Send recipient OTP
// @Description Send recipient OTP
// @Tags recipient
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Recipient ID"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /recipient/sendotp/{id} [post]
func (h *Handler) SendRecipientOTP(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	ID := c.Params("id")
	err := h.Service.SendRecipientOTP(c.Context(), ID, userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": OTPSent,
	})
}

// @Summary Activate recipient
// @Description Activate recipient
// @Tags recipient
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Recipient ID"
// @Param body body ActivateReq true "Activate request"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /recipient/activate/{id} [post]
func (h *Handler) ActivateRecipient(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	ID := c.Params("id")

	req := ActivateReq{}
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

	err = h.Service.ActivateRecipient(c.Context(), ID, userID, req.OTP)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": ActivateRecipientSuccess,
	})
}

// @Summary Delete recipient
// @Description Delete recipient
// @Tags recipient
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Recipient ID"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /recipient/{id} [delete]
func (h *Handler) DeleteRecipient(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	ID := c.Params("id")
	err := h.Service.DeleteRecipient(c.Context(), ID, userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": DeleteRecipientSuccess,
	})
}
