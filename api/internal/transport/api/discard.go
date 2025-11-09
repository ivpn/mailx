package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetDiscards = "Unable to retrieve discards for this user."
)

type DiscardService interface {
	GetDiscards(context.Context, string) ([]model.Discard, error)
}

// @Summary Get discards
// @Description Get all discards for the authenticated user
// @Tags discard
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} model.Discard
// @Failure 400 {object} ErrorRes
// @Router /discards [get]
func (h *Handler) GetDiscards(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	discards, err := h.Service.GetDiscards(c.Context(), userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrGetDiscards,
		})
	}

	return c.JSON(discards)
}
