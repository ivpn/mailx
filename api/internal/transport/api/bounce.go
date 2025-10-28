package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetBounces    = "Unable to retrieve bounces for this user."
	ErrGetBounceFile = "Unable to retrieve bounce file."
)

type BounceService interface {
	GetBounces(context.Context, string) ([]model.Bounce, error)
	GetBounceFile(context.Context, string, string) ([]byte, error)
}

// @Summary Get bounces
// @Description Get all bounces for the authenticated user
// @Tags bounce
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} model.Bounce
// @Failure 400 {object} ErrorRes
// @Router /bounces [get]
func (h *Handler) GetBounces(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	bounces, err := h.Service.GetBounces(c.Context(), userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrGetBounces,
		})
	}

	return c.JSON(bounces)
}

// @Summary Get bounce file
// @Description Get bounce file by ID for the authenticated user
// @Tags bounce
// @Accept json
// @Produce plain
// @Security ApiKeyAuth
// @Param id path string true "Bounce ID"
// @Success 200 {string} string "Bounce file content"
// @Failure 400 {object} ErrorRes
// @Router /bounce/file/{id} [get]
func (h *Handler) GetBounceFile(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	bounceID := c.Params("id")
	data, err := h.Service.GetBounceFile(c.Context(), bounceID, userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrGetBounceFile,
		})
	}

	c.Set("Content-Type", "text/plain")
	return c.Send(data)
}
