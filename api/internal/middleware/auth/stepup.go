package auth

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/model"
)

const StepUpRequiredCode = 80001

var ErrStepUpRequired = "Additional verification is required to continue."

type StepUpService interface {
	GetUser(context.Context, string) (model.User, error)
	ConsumeStepUp(context.Context, string) (bool, error)
}

// NewStepUp gates a route behind a fresh password/passkey confirmation. The grant is single-use:
// it is consumed by the first request that passes, so it must be re-verified for the next one.
func NewStepUp(service StepUpService) fiber.Handler {

	return func(c *fiber.Ctx) error {
		ok, err := service.ConsumeStepUp(c.Context(), GetAuthnCookie(c))
		if err == nil && ok {
			return c.Next()
		}

		user, err := service.GetUser(c.Context(), GetUserID(c))
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrStepUpRequired,
			"code":    StepUpRequiredCode,
			"methods": fiber.Map{
				"password": user.PasswordHash != "",
				"passkey":  len(user.Creds) > 0,
			},
		})
	}
}
