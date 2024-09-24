package api

import (
	"context"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/model"
)

var (
	ErrBeginRegistration = "could not begin registration"
)

type SessionService interface {
	GetSession(context.Context, string) (webauthn.SessionData, bool, error)
	SaveSession(context.Context, webauthn.SessionData, string, string) error
	DeleteSession(context.Context, string) error
}

func (h *Handler) BeginRegistration(c *fiber.Ctx) error {
	// Parse the request
	req := RecipientReq{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	// Validate the request
	err = h.Validator.Struct(req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	// Create new user
	user := model.User{
		Email:    req.Email,
		IsActive: false,
	}

	// Save the user
	err = h.Service.PostUser(c.Context(), user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Begin registration
	_, sessionData, err := h.WebAuthn.BeginRegistration(user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrBeginRegistration,
		})
	}

	// Save the session
	token := model.GenSessionToken()
	err = h.Service.SaveSession(c.Context(), *sessionData, token, user.ID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return nil
}

func (h *Handler) FinishRegistration(c *fiber.Ctx) error {
	return nil
}

func (h *Handler) BeginLogin(c *fiber.Ctx) error {
	return nil
}

func (h *Handler) FinishLogin(c *fiber.Ctx) error {
	return nil
}
