package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email-service/internal/model"
	"ivpn.net/email-service/internal/utils"
)

var (
	RegisterSuccess = "User created"
	LoginSuccess    = "Login successful"
	LogoutSuccess   = "Logout successful"
	ActivateSuccess = "User activated"
	CookieAuth      = "Auth"
)

type UserService interface {
	PostUser(context.Context, model.User) error
	GetUserByCredentials(context.Context, string, string) (model.User, error)
	ActivateUser(context.Context, string, string) error
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ActivateRequest struct {
	UserID string `json:"id"`
	OTP    string `json:"otp"`
}

func (h *Handler) Register(c *fiber.Ctx) error {
	// Parse the request
	req := UserRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Create new user
	user := model.User{
		Email:         req.Email,
		PasswordPlain: &req.Password,
	}

	// Save the user
	err = h.Service.PostUser(c.Context(), user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": RegisterSuccess,
	})
}

func (h *Handler) Login(c *fiber.Ctx) error {
	// Parse the request
	req := UserRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get the user
	user, err := h.Service.GetUserByCredentials(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Create auth token
	token, err := utils.CreateToken(h.Cfg, user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Set the token in ecrypted cookie
	c.Cookie(&fiber.Cookie{
		Name:  CookieAuth,
		Value: token,
	})

	return c.Status(200).JSON(fiber.Map{
		"message": LoginSuccess,
	})
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	c.ClearCookie(CookieAuth)

	return c.Status(200).JSON(fiber.Map{
		"message": LogoutSuccess,
	})
}

func (h *Handler) Activate(c *fiber.Ctx) error {
	req := ActivateRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = h.Service.ActivateUser(c.Context(), req.UserID, req.OTP)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": ActivateSuccess,
	})
}
