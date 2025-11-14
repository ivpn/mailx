package api

import (
	"context"
	"log"

	"github.com/araddon/dateparse"
	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetAccessKeys   = "Unable to retrieve access keys for this user."
	ErrPostAccessKey   = "Unable to create access key. Please try again."
	ErrDeleteAccessKey = "Unable to delete access key. Please try again."
)

type AccessKeyService interface {
	GetAccessKeys(context.Context, string) ([]model.AccessKey, error)
	PostAccessKey(context.Context, string, model.AccessKey) error
	DeleteAccessKey(context.Context, string, string) error
}

// @Summary Get access keys
// @Description Get all access keys for the authenticated user
// @Tags access_key
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} model.AccessKey
// @Failure 400 {object} ErrorRes
// @Router /access_keys [get]
func (h *Handler) GetAccessKeys(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	accessKeys, err := h.Service.GetAccessKeys(c.Context(), userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrGetAccessKeys,
		})
	}

	return c.JSON(accessKeys)
}

// @Summary Create access key
// @Description Create a new access key for the authenticated user
// @Tags access_key
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param access_key body AccessKeyReq true "Access Key Request"
// @Success 200 {object} map[string]string "access_key"
// @Failure 400 {object} ErrorRes
// @Router /access_key [post]
func (h *Handler) PostAccessKey(c *fiber.Ctx) error {
	// Parse the request
	userId := auth.GetUserID(c)
	req := AccessKeyReq{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Printf("error login: %s", err.Error())
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	// Validate the request
	err = h.Validator.Struct(req)
	if err != nil {
		log.Printf("error login: %s", err.Error())
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	// Create token
	token, err := model.GenSessionToken()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrPostAccessKey,
		})
	}

	// Prepare access key model
	accessKey := model.AccessKey{
		UserID:     userId,
		TokenPlain: &token,
		ExpiresAt:  model.NeverExpires(),
	}

	// Set expiration if provided
	if req.ExpiresAt != "" {
		expiresAt, err := dateparse.ParseAny(req.ExpiresAt)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": ErrPostAccessKey,
			})
		}
		accessKey.ExpiresAt = &expiresAt
	}

	// Store access key
	err = h.Service.PostAccessKey(c.Context(), userId, accessKey)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrPostAccessKey,
		})
	}

	// Return the token
	return c.Status(200).JSON(fiber.Map{
		"access_key": token,
	})
}

// @Summary Delete access key
// @Description Delete access key by ID
// @Tags access_key
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Access Key ID"
// @Success 200
// @Failure 400 {object} ErrorRes
// @Router /access_key/{id} [delete]
func (h *Handler) DeleteAccessKey(c *fiber.Ctx) error {
	userId := auth.GetUserID(c)
	id := c.Params("id")
	err := h.Service.DeleteAccessKey(c.Context(), id, userId)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrDeleteAccessKey,
		})
	}

	return nil
}
