package api

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email-service/internal/middleware/auth"
	"ivpn.net/email-service/internal/model"
	"ivpn.net/email-service/internal/utils"
)

var (
	RegisterSuccess       = "User created"
	LoginSuccess          = "Login successful"
	LogoutSuccess         = "Logout successful"
	DeleteUserSuccess     = "User deleted"
	OTPSent               = "OTP sent"
	ActivateUserSuccess   = "User activated"
	ErrInvalidCredentials = "Invalid credentials"
	ErrInvalidRequest     = "Invalid request"
)

type UserService interface {
	PostUser(context.Context, model.User) error
	SendUserOTP(context.Context, string) error
	ActivateUser(context.Context, string, string) error
	GetUserByCredentials(context.Context, string, string) (model.User, error)
	GetUserByPassword(context.Context, string, string) (model.User, error)
	DeleteUser(context.Context, string) error
}

func (h *Handler) Register(c *fiber.Ctx) error {
	// Parse the request
	req := UserReq{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	// Create new user
	user := model.User{
		Email:         req.Email,
		PasswordPlain: &req.Password,
		IsActive:      false,
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

func (h *Handler) SendUserOTP(c *fiber.Ctx) error {
	ID := auth.GetUserID(c)

	err := h.Service.SendUserOTP(c.Context(), ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": OTPSent,
	})
}

func (h *Handler) Activate(c *fiber.Ctx) error {
	ID := auth.GetUserID(c)

	req := ActivateReq{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	err = h.Service.ActivateUser(c.Context(), ID, req.OTP)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": ActivateUserSuccess,
	})
}

func (h *Handler) Login(c *fiber.Ctx) error {
	// Parse the request
	req := UserReq{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Printf("error login: %s", err.Error())
		return c.Status(500).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	// Validate the request
	userModel := model.User{
		Email:         req.Email,
		PasswordPlain: &req.Password,
	}
	err = userModel.Validate()
	if err != nil {
		log.Printf("error login: %s", err.Error())
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidCredentials,
		})
	}

	// Get the user
	user, err := h.Service.GetUserByCredentials(c.Context(), req.Email, req.Password)
	if err != nil {
		log.Printf("error login: %s", err.Error())
		return c.Status(401).JSON(fiber.Map{
			"error": ErrInvalidCredentials,
		})
	}

	// Create auth token
	token, err := utils.CreateAuthToken(h.Cfg, user.ID)
	if err != nil {
		log.Printf("error login: %s", err.Error())
		return c.Status(500).JSON(fiber.Map{
			"error": ErrInvalidCredentials,
		})
	}

	// Set the token in encrypted cookie
	c.Cookie(&fiber.Cookie{
		Name:  auth.AUTH_COOKIE,
		Value: token,
	})

	return c.Status(200).JSON(fiber.Map{
		"message": LoginSuccess,
	})
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	c.ClearCookie(auth.AUTH_COOKIE)

	return c.Status(200).JSON(fiber.Map{
		"message": LogoutSuccess,
	})
}

func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	ID := auth.GetUserID(c)

	// Parse the request
	req := DeleteUserReq{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Printf("error deleting user: %s", err.Error())
		return c.Status(500).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	// Validate the request
	err = utils.ValidatePassword(req.Password)
	if err != nil {
		log.Printf("error deleting user: %s", err.Error())
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidCredentials,
		})
	}

	// Get the user
	user, err := h.Service.GetUserByPassword(c.Context(), ID, req.Password)
	if err != nil || user.ID != ID {
		log.Printf("error deleting user: %s", err.Error())
		return c.Status(401).JSON(fiber.Map{
			"error": ErrInvalidCredentials,
		})
	}

	// Delete the user
	err = h.Service.DeleteUser(c.Context(), ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": DeleteUserSuccess,
	})
}
