package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetLogs        = "Unable to retrieve logs for this user."
	ErrGetLogFile     = "Unable to retrieve log file."
	ErrDeleteLogs     = "Unable to delete logs for this user."
	LogsDeleteSuccess = "All diagnostic logs have been deleted."
)

type LogService interface {
	GetLogs(context.Context, string) ([]model.Log, error)
	GetLogFile(context.Context, string, string) ([]byte, error)
	DeleteLogs(context.Context, string) error
}

// @Summary Get logs
// @Description Get all logs for the authenticated user
// @Tags log
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} model.Log
// @Failure 400 {object} ErrorRes
// @Router /logs [get]
func (h *Handler) GetLogs(c *fiber.Ctx) error {
	userId := auth.GetUserID(c)
	logs, err := h.Service.GetLogs(c.Context(), userId)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrGetLogs,
		})
	}

	return c.JSON(logs)
}

// @Summary Get log file
// @Description Get log file by ID for the authenticated user
// @Tags log
// @Accept json
// @Produce plain
// @Security ApiKeyAuth
// @Param id path string true "Log ID"
// @Success 200 {string} string "Log file content"
// @Failure 400 {object} ErrorRes
// @Router /log/file/{id} [get]
func (h *Handler) GetLogFile(c *fiber.Ctx) error {
	userId := auth.GetUserID(c)
	logId := c.Params("id")
	data, err := h.Service.GetLogFile(c.Context(), logId, userId)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrGetLogFile,
		})
	}

	c.Set("Content-Type", "text/plain")
	return c.Send(data)
}

// @Summary Delete logs
// @Description Delete all logs for the authenticated user
// @Tags log
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /logs [delete]
func (h *Handler) DeleteLogs(c *fiber.Ctx) error {
	userId := auth.GetUserID(c)
	err := h.Service.DeleteLogs(c.Context(), userId)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrDeleteLogs,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": LogsDeleteSuccess,
	})
}
