package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetDomains   = "Unable to retrieve custom domains for this user."
	ErrGetDNSConfig = "Unable to retrieve custom domains DNS config for this user."
	ErrPostDomain   = "Unable to create custom domain. Please try again."
	ErrUpdateDomain = "Unable to update custom domain. Please try again."
	ErrDeleteDomain = "Unable to delete custom domain. Please try again."
)

type DomainService interface {
	GetDomains(context.Context, string) ([]model.Domain, error)
	GetDNSConfig(context.Context, string) (model.DNSConfig, error)
	PostDomain(context.Context, model.Domain) (model.Domain, error)
	UpdateDomain(context.Context, model.Domain) error
	DeleteDomain(context.Context, string, string) error
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
