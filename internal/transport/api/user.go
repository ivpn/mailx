package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email-service/internal/model"
	"ivpn.net/email-service/internal/utils"
)

var (
	LoginSuccess = "Login successful"
)

type UserService interface {
	PostUser(context.Context, model.User) error
	Login(context.Context, string, string) (model.User, error)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(c *fiber.Ctx) error {
	req := LoginRequest{}

	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.Service.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	token, err := utils.CreateToken(h.Cfg, user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:  "Auth",
		Value: token,
	})

	return c.Status(200).JSON(fiber.Map{
		"message": LoginSuccess,
	})
}
