package api

import (
	"context"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
)

var (
	PostAliasSuccess   = "Alias created successfully."
	UpdateAliasSuccess = "Alias updated successfully."
	DeleteAliasSuccess = "Alias deleted successfully."
	ErrInvalidDomain   = "Selected domain is invalid."
	ErrUnverifiedRcp   = "The recipient address has not been verified."
)

type AliasService interface {
	GetAlias(context.Context, string, string) (model.Alias, error)
	GetAliases(context.Context, string, int, int, string, string, string, string) (model.AliasList, error)
	GetAllAliases(context.Context, string) ([]model.Alias, error)
	PostAlias(context.Context, model.Alias, string, string, string) (model.Alias, error)
	UpdateAlias(context.Context, model.Alias) error
	DeleteAlias(context.Context, string, string) error
}

// @Summary Get alias
// @Description Get alias by ID
// @Tags alias
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Alias ID"
// @Success 200 {object} model.Alias
// @Failure 400 {object} ErrorRes
// @Router /alias/{id} [get]
func (h *Handler) GetAlias(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	id := c.Params("id")
	alias, err := h.Service.GetAlias(c.Context(), id, userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(alias)
}

// @Summary Get aliases
// @Description Get all aliases
// @Tags alias
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} model.AliasList
// @Failure 400 {object} ErrorRes
// @Router /aliases [get]
func (h *Handler) GetAliases(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 0
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 0
	}

	sortBy := c.Query("sort_by")
	sortOrder := strings.ToUpper(c.Query("sort_order"))
	catchAll := c.Query("catch_all")
	search := c.Query("search")

	var allowSortBy = map[string]bool{
		"created_at": true,
		"updated_at": true,
		"name":       true,
	}
	var allowSortOrder = map[string]bool{
		"ASC":  true,
		"DESC": true,
	}
	var allowCatchAll = map[string]bool{
		"true":  true,
		"false": true,
		"":      true,
	}

	if _, ok := allowSortBy[sortBy]; !ok {
		sortBy = "created_at"
	}
	if _, ok := allowSortOrder[sortOrder]; !ok {
		sortOrder = "DESC"
	}
	if _, ok := allowCatchAll[catchAll]; !ok {
		catchAll = ""
	}

	err = h.Validator.Var(search, "omitempty,required,search")
	if err != nil {
		return c.JSON(model.AliasList{
			Aliases: []model.Alias{},
			Total:   0,
		})
	}

	list, err := h.Service.GetAliases(c.Context(), userID, limit, page, sortBy, sortOrder, catchAll, search)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(list)
}

// @Summary Export aliases
// @Description Export all aliases as CSV
// @Tags alias
// @Accept json
// @Produce text/csv
// @Security ApiKeyAuth
// @Success 200 {string} string "CSV data"
// @Failure 400 {object} ErrorRes
// @Router /aliases/export [get]
func (h *Handler) ExportAliases(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	aliases, err := h.Service.GetAllAliases(c.Context(), userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=\"aliases.csv\"")

	var b strings.Builder
	b.WriteString("alias,description,enabled,recipients\n")
	for _, alias := range aliases {
		b.WriteString(alias.Name + ",")
		b.WriteString(alias.Description + ",")
		b.WriteString(strconv.FormatBool(alias.Enabled) + ",")
		b.WriteString(alias.Recipients + "\n")
	}

	return c.SendString(b.String())
}

// @Summary Create alias
// @Description Create alias
// @Tags alias
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body AliasReq true "Alias request"
// @Success 201 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /alias [post]
func (h *Handler) PostAlias(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	req := AliasReq{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	err = h.Validator.Struct(req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	if !strings.Contains(h.Cfg.Domains, req.Domain) {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidDomain,
		})
	}

	rcps, err := h.Service.GetVerifiedRecipients(c.Context(), req.Recipients, userID)
	if err != nil || len(rcps) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrUnverifiedRcp,
		})
	}

	if req.Format == model.AliasFormatCatchAll && req.CatchAllSuffix == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	alias := model.Alias{
		UserID:      userID,
		Description: req.Description,
		Enabled:     req.Enabled,
		Recipients:  model.GetEmails(rcps),
		FromName:    req.FromName,
	}

	alias, err = h.Service.PostAlias(c.Context(), alias, req.Format, req.Domain, req.CatchAllSuffix)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": PostAliasSuccess,
		"alias":   alias,
	})
}

// @Summary Update alias
// @Description Update alias
// @Tags alias
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Alias ID"
// @Param body body AliasReq true "Alias request"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /alias/{id} [put]
func (h *Handler) UpdateAlias(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	req := AliasReq{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	rcps, err := h.Service.GetVerifiedRecipients(c.Context(), req.Recipients, userID)
	if err != nil || len(rcps) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	alias := model.Alias{
		UserID:      userID,
		Description: req.Description,
		Enabled:     req.Enabled,
		Recipients:  model.GetEmails(rcps),
		FromName:    req.FromName,
	}
	alias.ID = c.Params("id")

	err = h.Service.UpdateAlias(c.Context(), alias)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": UpdateAliasSuccess,
	})
}

// @Summary Delete alias
// @Description Delete alias
// @Tags alias
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Alias ID"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /alias/{id} [delete]
func (h *Handler) DeleteAlias(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	id := c.Params("id")
	err := h.Service.DeleteAlias(c.Context(), id, userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": DeleteAliasSuccess,
	})
}
