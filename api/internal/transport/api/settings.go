package api

import (
	"context"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
)

var (
	UpdateSettingsSuccess      = "Settings updated successfully."
	UpdateBlockTrackersSuccess = "Block email trackers setting updated successfully."
	ErrBlockTrackersRequired   = "The block_trackers field is required."
	ErrBlockTrackersInvalid    = "The block_trackers field must be true or false."
	UpdateRemoveHeaderSuccess  = "Mailx header setting updated successfully."
	ErrRemoveHeaderRequired    = "The remove_header field is required."
	ErrRemoveHeaderInvalid     = "The remove_header field must be true or false."
	UpdateLogIssuesSuccess     = "Diagnostics setting updated successfully."
	ErrLogIssuesRequired       = "The log_issues field is required."
	ErrLogIssuesInvalid        = "The log_issues field must be true or false."
	UpdateDefaultDomainSuccess = "Default domain updated successfully."
	ErrDefaultDomainRequired   = "The domain field is required."
	ErrDefaultDomainInvalid    = "The selected domain is invalid."
	UpdateDefaultRcpSuccess    = "Default recipient updated successfully."
	ErrDefaultRcpRequired      = "The recipient field is required."
	ErrDefaultRcpInvalid       = "The selected recipient is invalid."
	UpdateAliasFormatSuccess   = "Default alias format updated successfully."
	ErrAliasFormatRequired     = "The alias_format field is required."
	ErrAliasFormatInvalid      = "The selected alias format is invalid."
	UpdateFromNameSuccess      = "From name updated successfully."
	ErrFromNameRequired        = "The from_name field is required."
	ErrFromNameInvalid         = "The from name must be at most 64 characters and can't contain line breaks."
)

// MaxFromNameLength is the longest From name, in characters.
const MaxFromNameLength = 64

// defaultAliasFormats are the formats a new alias can be created with by default.
var defaultAliasFormats = []string{model.AliasFormatRandomWords, model.AliasFormatRandomChars, model.AliasFormatUUID}

type SettingsService interface {
	GetSettings(context.Context, string) (model.Settings, error)
	UpdateSettings(context.Context, model.Settings) error
	UpdateBlockTrackers(context.Context, string, bool) error
	UpdateRemoveHeader(context.Context, string, bool) error
	UpdateLogIssues(context.Context, string, bool) error
	UpdateDefaultDomain(context.Context, string, string) error
	UpdateDefaultRecipient(context.Context, string, string) error
	UpdateAliasFormat(context.Context, string, string) error
	UpdateFromName(context.Context, string, string) error
}

// @Summary Get settings
// @Description Get settings
// @Tags settings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} model.Settings
// @Failure 400 {object} ErrorRes
// @Router /settings [get]
func (h *Handler) GetSettings(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)

	settings, err := h.Service.GetSettings(c.Context(), userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(settings)
}

// @Summary Update settings
// @Description Update settings
// @Tags settings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param settings body SettingsReq true "Settings"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /settings [put]
func (h *Handler) UpdateSettings(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)

	req := SettingsReq{}
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

	settings := model.Settings{
		UserID:       userID,
		Domain:       req.Domain,
		Recipient:    req.Recipient,
		FromName:     req.FromName,
		AliasFormat:  req.AliasFormat,
		LogIssues:    req.LogIssues,
		RemoveHeader: req.RemoveHeader,
	}
	settings.ID = req.ID

	err = h.Service.UpdateSettings(c.Context(), settings)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": UpdateSettingsSuccess,
	})
}

// @Summary Update block email trackers
// @Description Turn removal of tracking pixels from forwarded emails on or off
// @Tags settings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param settings body BlockTrackersReq true "Block email trackers"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Failure 422 {object} ErrorRes
// @Router /settings/block-trackers [put]
func (h *Handler) UpdateBlockTrackers(c *fiber.Ctx) error {
	req := BlockTrackersReq{}
	return h.updateBoolSetting(c, &req, func() *bool { return req.BlockTrackers }, boolSetting{
		errRequired: ErrBlockTrackersRequired,
		errInvalid:  ErrBlockTrackersInvalid,
		success:     UpdateBlockTrackersSuccess,
		update:      h.Service.UpdateBlockTrackers,
	})
}

// @Summary Update Mailx header
// @Description Turn the Mailx header in forwarded emails off (remove_header true) or on
// @Tags settings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param settings body RemoveHeaderReq true "Remove Mailx header"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Failure 422 {object} ErrorRes
// @Router /settings/remove-header [put]
func (h *Handler) UpdateRemoveHeader(c *fiber.Ctx) error {
	req := RemoveHeaderReq{}
	return h.updateBoolSetting(c, &req, func() *bool { return req.RemoveHeader }, boolSetting{
		errRequired: ErrRemoveHeaderRequired,
		errInvalid:  ErrRemoveHeaderInvalid,
		success:     UpdateRemoveHeaderSuccess,
		update:      h.Service.UpdateRemoveHeader,
	})
}

// @Summary Update diagnostics
// @Description Turn diagnostic logging of failed deliveries and forwarding issues on or off
// @Tags settings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param settings body LogIssuesReq true "Enable diagnostics"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Failure 422 {object} ErrorRes
// @Router /settings/log-issues [put]
func (h *Handler) UpdateLogIssues(c *fiber.Ctx) error {
	req := LogIssuesReq{}
	return h.updateBoolSetting(c, &req, func() *bool { return req.LogIssues }, boolSetting{
		errRequired: ErrLogIssuesRequired,
		errInvalid:  ErrLogIssuesInvalid,
		success:     UpdateLogIssuesSuccess,
		update:      h.Service.UpdateLogIssues,
	})
}

// boolSetting describes a single on/off setting updated through its own endpoint.
type boolSetting struct {
	errRequired string
	errInvalid  string
	success     string
	update      func(context.Context, string, bool) error
}

// updateBoolSetting parses req from the body and stores the value returned by value, which
// must be a required JSON boolean.
func (h *Handler) updateBoolSetting(c *fiber.Ctx, req any, value func() *bool, s boolSetting) error {
	userID := auth.GetUserID(c)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": s.errInvalid,
		})
	}
	v := value()
	if v == nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": s.errRequired,
		})
	}

	err := s.update(c.Context(), userID, *v)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": s.success,
	})
}

// @Summary Update default domain
// @Description Set the default alias domain for new aliases: a built-in domain or a verified custom domain name
// @Tags settings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param settings body DefaultDomainReq true "Default domain"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Failure 422 {object} ErrorRes
// @Router /settings/domain [put]
func (h *Handler) UpdateDefaultDomain(c *fiber.Ctx) error {
	req := DefaultDomainReq{}
	return h.updateStringSetting(c, &req, func() *string { return req.Domain }, stringSetting{
		errRequired: ErrDefaultDomainRequired,
		errInvalid:  ErrDefaultDomainInvalid,
		success:     UpdateDefaultDomainSuccess,
		valid:       h.isUserDomain,
		update:      h.Service.UpdateDefaultDomain,
	})
}

// @Summary Update default recipient
// @Description Set the default recipient for new aliases; it must be one of the user's verified recipients
// @Tags settings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param settings body DefaultRecipientReq true "Default recipient"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Failure 422 {object} ErrorRes
// @Router /settings/recipient [put]
func (h *Handler) UpdateDefaultRecipient(c *fiber.Ctx) error {
	req := DefaultRecipientReq{}
	return h.updateStringSetting(c, &req, func() *string { return req.Recipient }, stringSetting{
		errRequired: ErrDefaultRcpRequired,
		errInvalid:  ErrDefaultRcpInvalid,
		success:     UpdateDefaultRcpSuccess,
		valid:       h.isUserRecipient,
		update:      h.Service.UpdateDefaultRecipient,
	})
}

// @Summary Update default alias format
// @Description Set the default alias format for new aliases: words, random or uuid
// @Tags settings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param settings body AliasFormatReq true "Default alias format"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Failure 422 {object} ErrorRes
// @Router /settings/alias-format [put]
func (h *Handler) UpdateAliasFormat(c *fiber.Ctx) error {
	req := AliasFormatReq{}
	return h.updateStringSetting(c, &req, func() *string { return req.AliasFormat }, stringSetting{
		errRequired: ErrAliasFormatRequired,
		errInvalid:  ErrAliasFormatInvalid,
		success:     UpdateAliasFormatSuccess,
		valid: func(ctx context.Context, userID string, format string) bool {
			return slices.Contains(defaultAliasFormats, format)
		},
		update: h.Service.UpdateAliasFormat,
	})
}

// @Summary Update from name
// @Description Set the From name used for replies and sends; an empty value uses the alias address
// @Tags settings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param settings body FromNameReq true "From name"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Failure 422 {object} ErrorRes
// @Router /settings/from-name [put]
func (h *Handler) UpdateFromName(c *fiber.Ctx) error {
	req := FromNameReq{}
	return h.updateStringSetting(c, &req, func() *string { return req.FromName }, stringSetting{
		errRequired: ErrFromNameRequired,
		errInvalid:  ErrFromNameInvalid,
		success:     UpdateFromNameSuccess,
		allowEmpty:  true,
		valid: func(ctx context.Context, userID string, name string) bool {
			// Control characters could break out of the From header
			return utf8.RuneCountInString(name) <= MaxFromNameLength && !strings.ContainsFunc(name, unicode.IsControl)
		},
		update: h.Service.UpdateFromName,
	})
}

// isUserDomain reports whether domain is a built-in domain or one of the user's verified
// custom domains.
func (h *Handler) isUserDomain(ctx context.Context, userID string, domain string) bool {
	if slices.Contains(strings.Split(h.Cfg.Domains, ","), domain) {
		return true
	}

	domains, err := h.Service.GetVerifiedDomains(ctx, userID)
	if err != nil {
		return false
	}

	return slices.ContainsFunc(domains, func(d model.Domain) bool { return d.Name == domain })
}

// isUserRecipient reports whether email is one of the user's verified recipients.
func (h *Handler) isUserRecipient(ctx context.Context, userID string, email string) bool {
	// The lookup takes a comma-separated list, so only accept a single address
	if strings.Contains(email, ",") {
		return false
	}

	rcps, err := h.Service.GetVerifiedRecipients(ctx, email, userID)
	return err == nil && len(rcps) == 1
}

// stringSetting describes a single text setting updated through its own endpoint.
type stringSetting struct {
	errRequired string
	errInvalid  string
	success     string
	allowEmpty  bool
	valid       func(context.Context, string, string) bool
	update      func(context.Context, string, string) error
}

// updateStringSetting parses req from the body and stores the trimmed value returned by value,
// which must be a JSON string that passes s.valid.
func (h *Handler) updateStringSetting(c *fiber.Ctx, req any, value func() *string, s stringSetting) error {
	userID := auth.GetUserID(c)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": s.errInvalid,
		})
	}
	v := value()
	if v == nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": s.errRequired,
		})
	}
	val := strings.TrimSpace(*v)
	if val == "" && !s.allowEmpty {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": s.errRequired,
		})
	}
	if val != "" && !s.valid(c.Context(), userID, val) {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": s.errInvalid,
		})
	}

	err := s.update(c.Context(), userID, val)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": s.success,
	})
}
