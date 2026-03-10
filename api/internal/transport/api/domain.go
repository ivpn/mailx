package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetDomains                = "Unable to retrieve custom domains for this user."
	ErrGetDomain                 = "Unable to retrieve custom domain for this user."
	ErrGetDNSConfig              = "Unable to retrieve custom domains DNS config for this user."
	ErrPostDomain                = "Unable to create custom domain. Please try again."
	ErrUpdateDomain              = "Unable to update custom domain. Please try again."
	ErrDeleteDomain              = "Unable to delete custom domain. Please try again."
	PostDomainSuccess            = "Custom domain added successfully."
	UpdateDomainSuccess          = "Custom domain updated successfully."
	DeleteDomainSuccess          = "Custom domain deleted successfully."
	DNSRecordVerificationSuccess = "Custom domain DNS records verified successfully."
)

type DomainService interface {
	GetDomains(context.Context, string) ([]model.Domain, error)
	GetDomain(context.Context, string, string) (model.Domain, error)
	GetDNSConfig(context.Context, string) (model.DNSConfig, error)
	PostDomain(context.Context, model.Domain) (model.Domain, error)
	UpdateDomain(context.Context, model.Domain) error
	DeleteDomain(context.Context, string, string) error
	VerifyDomainDNSRecords(context.Context, string, string) error
}

// @Summary Get custom domains
// @Description Get all custom domains for the authenticated user
// @Tags domain
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} model.Domain
// @Failure 400 {object} ErrorRes
// @Router /domains [get]
func (h *Handler) GetDomains(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	domains, err := h.Service.GetDomains(c.Context(), userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrGetDomains,
		})
	}

	return c.JSON(domains)
}

// @Summary Get custom domains DNS config
// @Description Get the DNS configuration for all custom domains of the authenticated user
// @Tags domain
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} model.DNSConfig
// @Failure 400 {object} ErrorRes
// @Router /domains/dns-config [get]
func (h *Handler) GetDNSConfig(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	dnsConfig, err := h.Service.GetDNSConfig(c.Context(), userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrGetDNSConfig,
		})
	}

	return c.JSON(dnsConfig)
}

// @Summary Create custom domain
// @Description Create a new custom domain for the authenticated user
// @Tags domain
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param domain body DomainReq true "Custom Domain Request"
// @Success 201 {object} map[string]string "message"
// @Failure 400 {object} ErrorRes
// @Router /domains [post]
func (h *Handler) PostDomain(c *fiber.Ctx) error {
	// Parse the request
	userID := auth.GetUserID(c)
	req := DomainReq{}
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

	// Create domain
	domain := model.Domain{
		UserID:  userID,
		Name:    req.Name,
		Enabled: true,
	}

	// Post domain
	_, err = h.Service.PostDomain(c.Context(), domain)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": PostDomainSuccess,
	})
}

// @Summary Update custom domain
// @Description Update an existing custom domain for the authenticated user
// @Tags domain
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param domain body UpdateDomainReq true "Update Custom Domain Request"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} ErrorRes
// @Router /domains [put]
func (h *Handler) UpdateDomain(c *fiber.Ctx) error {
	// Parse the request
	userID := auth.GetUserID(c)
	req := UpdateDomainReq{}
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

	// Get existing domain
	domain, err := h.Service.GetDomain(c.Context(), req.ID, userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrGetDomain,
		})
	}

	// Update domain fields
	domain.Description = req.Description
	domain.Recipient = req.Recipient
	domain.FromName = req.FromName
	domain.Enabled = req.Enabled

	// Update domain
	err = h.Service.UpdateDomain(c.Context(), domain)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrUpdateDomain,
		})
	}

	return c.JSON(fiber.Map{
		"message": UpdateDomainSuccess,
	})
}

// @Summary Delete custom domain
// @Description Delete an existing custom domain for the authenticated user
// @Tags domain
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Domain ID"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} ErrorRes
// @Router /domains/{id} [delete]
func (h *Handler) DeleteDomain(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	domainID := c.Params("id")

	err := h.Service.DeleteDomain(c.Context(), domainID, userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrDeleteDomain,
		})
	}

	return c.JSON(fiber.Map{
		"message": DeleteDomainSuccess,
	})
}

// @Summary Verify custom domain DNS records
// @Description Verify the DNS records for a custom domain of the authenticated user
// @Tags domain
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Domain ID"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} ErrorRes
// @Router /domains/{id}/verify-dns [post]
func (h *Handler) VerifyDomainDNSRecords(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	domainID := c.Params("id")

	err := h.Service.VerifyDomainDNSRecords(c.Context(), domainID, userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": DNSRecordVerificationSuccess,
	})
}
