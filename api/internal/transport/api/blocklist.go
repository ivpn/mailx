package api

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetBlocklist            = "Unable to retrieve blocklist for this user."
	ErrInvalidBlocklistEmail   = "Please enter a valid email address."
	ErrInvalidBlocklistDomain  = "Please enter a valid domain."
	PostBlocklistItemSuccess   = "Added to blocklist."
	DeleteBlocklistItemSuccess = "Removed from blocklist."
)

// maxBlocklistLimit caps the page size, same as maxAliasLimit does for aliases.
const maxBlocklistLimit = 100

type BlocklistService interface {
	GetBlocklist(context.Context, string, int, int) (model.BlocklistList, error)
	PostBlocklistItem(context.Context, model.BlocklistItem) (model.BlocklistItem, error)
	DeleteBlocklistItem(context.Context, string, string) error
}

// @Summary Get blocklist
// @Description Get a page of blocklist entries for the authenticated user
// @Tags blocklist
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param limit query int false "Page size (max 100)"
// @Param page query int false "Page number"
// @Success 200 {object} model.BlocklistList
// @Failure 400 {object} ErrorRes
// @Router /blocklist [get]
func (h *Handler) GetBlocklist(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 || limit > maxBlocklistLimit {
		limit = maxBlocklistLimit
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 0
	}

	list, err := h.Service.GetBlocklist(c.Context(), userID, limit, page)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrGetBlocklist,
		})
	}

	return c.JSON(list)
}

// @Summary Add to blocklist
// @Description Block inbound email from a sender address or domain (including its subdomains)
// @Tags blocklist
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param blocklist body BlocklistReq true "Blocklist Request"
// @Success 201 {object} map[string]string "id, message"
// @Failure 400 {object} ErrorRes
// @Router /blocklist [post]
func (h *Handler) PostBlocklistItem(c *fiber.Ctx) error {
	// Parse the request
	userID := auth.GetUserID(c)
	req := BlocklistReq{}
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

	// Validate the value against its type, in the form it's stored and matched in
	itemType := model.BlocklistType(req.Type)
	value := model.NormalizeBlocklistValue(itemType, req.Value)
	if err := h.Validator.Var(value, model.BlocklistValueRules(itemType)); err != nil {
		errMsg := ErrInvalidBlocklistDomain
		if itemType == model.BlocklistEmail {
			errMsg = ErrInvalidBlocklistEmail
		}
		return c.Status(400).JSON(fiber.Map{
			"error": errMsg,
		})
	}

	// Post blocklist item
	created, err := h.Service.PostBlocklistItem(c.Context(), model.BlocklistItem{
		UserID: userID,
		Type:   itemType,
		Value:  value,
	})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"id":      created.ID,
		"message": PostBlocklistItemSuccess,
	})
}

// @Summary Delete blocklist entry
// @Description Remove an entry from the authenticated user's blocklist
// @Tags blocklist
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Blocklist entry ID"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /blocklist/{id} [delete]
func (h *Handler) DeleteBlocklistItem(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	ID := c.Params("id")

	err := h.Service.DeleteBlocklistItem(c.Context(), ID, userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": DeleteBlocklistItemSuccess,
	})
}
