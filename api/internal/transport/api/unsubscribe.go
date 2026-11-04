package api

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils/signedurl"
)

var (
	ErrListUnsubscribeBehaviourRequired = "The behaviour field is required."
	ErrListUnsubscribeBehaviourInvalid  = "The selected behaviour is invalid."
	UpdateListUnsubscribeSuccess        = "Unsubscribe behaviour updated successfully."
)

type OneClickService interface {
	UpdateListUnsubscribeBehaviour(context.Context, string, model.ListUnsubscribeBehaviour) error
	OneClickSigner() *signedurl.Signer
	OneClickDeactivateAlias(context.Context, string, model.OneClickRequest) error
	OneClickDeleteAlias(context.Context, string, model.OneClickRequest) error
	OneClickBlock(context.Context, string, model.BlocklistType, string) error
}

// @Summary Update List-Unsubscribe behaviour
// @Description Set what a mail client's "Unsubscribe" button does on forwarded emails
// @Tags settings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param settings body ListUnsubscribeBehaviourReq true "List-Unsubscribe behaviour"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Failure 422 {object} ErrorRes
// @Router /settings/list-unsubscribe-behaviour [put]
func (h *Handler) UpdateListUnsubscribeBehaviour(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)

	req := ListUnsubscribeBehaviourReq{}
	err := c.BodyParser(&req)
	if err != nil || req.Behaviour == nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": ErrListUnsubscribeBehaviourRequired,
		})
	}

	behaviour := model.ListUnsubscribeBehaviour(*req.Behaviour)
	if !behaviour.Valid() {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": ErrListUnsubscribeBehaviourInvalid,
		})
	}

	err = h.Service.UpdateListUnsubscribeBehaviour(c.Context(), userID, behaviour)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": UpdateListUnsubscribeSuccess,
	})
}

// DeactivateOneClick handles an RFC 8058 one-click unsubscribe by deactivating the alias.
func (h *Handler) DeactivateOneClick(c *fiber.Ctx) error {
	err := h.Service.OneClickDeactivateAlias(c.Context(), c.Params("alias"), oneClickRequest(c))
	return oneClickResponse(c, err)
}

// DeleteOneClick handles an RFC 8058 one-click unsubscribe by deleting the alias.
func (h *Handler) DeleteOneClick(c *fiber.Ctx) error {
	err := h.Service.OneClickDeleteAlias(c.Context(), c.Params("alias"), oneClickRequest(c))
	return oneClickResponse(c, err)
}

// BlockEmailOneClick handles an RFC 8058 one-click unsubscribe by blocklisting the sender address.
func (h *Handler) BlockEmailOneClick(c *fiber.Ctx) error {
	err := h.Service.OneClickBlock(c.Context(), c.Params("alias"), model.BlocklistEmail, c.Query("email"))
	return oneClickResponse(c, err)
}

// BlockDomainOneClick handles an RFC 8058 one-click unsubscribe by blocklisting the sender domain.
func (h *Handler) BlockDomainOneClick(c *fiber.Ctx) error {
	err := h.Service.OneClickBlock(c.Context(), c.Params("alias"), model.BlocklistDomain, c.Query("domain"))
	return oneClickResponse(c, err)
}

// MethodNotAllowedOneClick keeps GETs from link scanners and prefetchers, and any other
// non-POST method, away from the one-click actions, which only run on POST.
func (h *Handler) MethodNotAllowedOneClick(c *fiber.Ctx) error {
	c.Set(fiber.HeaderAllow, fiber.MethodPost)
	c.Status(fiber.StatusMethodNotAllowed)
	return nil
}

// TooManyRequestsOneClick answers a rate-limited one-click request with an empty body.
func (h *Handler) TooManyRequestsOneClick(c *fiber.Ctx) error {
	c.Status(fiber.StatusTooManyRequests)
	return nil
}

func oneClickRequest(c *fiber.Ctx) model.OneClickRequest {
	return model.OneClickRequest{
		IP:        c.IP(),
		UserAgent: c.Get(fiber.HeaderUserAgent),
	}
}

// oneClickResponse answers with an empty body, so the response doesn't reveal the alias state.
func oneClickResponse(c *fiber.Ctx, err error) error {
	switch {
	case err == nil:
		c.Status(fiber.StatusOK)
		return nil
	case errors.Is(err, model.ErrOneClickAliasNotFound):
		c.Status(fiber.StatusNotFound)
	case errors.Is(err, model.ErrOneClickInvalidValue):
		c.Status(fiber.StatusUnprocessableEntity)
	default:
		c.Status(fiber.StatusInternalServerError)
	}

	return nil
}
