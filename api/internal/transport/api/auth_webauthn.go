package api

import (
	"context"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gofiber/fiber/v2"
)

type SessionService interface {
	GetSession(context.Context, string) (webauthn.SessionData, bool, error)
	SaveSession(context.Context, webauthn.SessionData, string) error
	DeleteSession(context.Context, string) error
}

func (h *Handler) BeginRegistration(ctx *fiber.Ctx) error {
	return nil
}

func (h *Handler) FinishRegistration(ctx *fiber.Ctx) error {
	return nil
}

func (h *Handler) BeginLogin(ctx *fiber.Ctx) error {
	return nil
}

func (h *Handler) FinishLogin(ctx *fiber.Ctx) error {
	return nil
}
