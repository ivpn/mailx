package api

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
)

var (
	PostRuleSuccess       = "Rule created successfully."
	UpdateRuleSuccess     = "Rule updated successfully."
	UpdateRuleActiveOn    = "Rule activated."
	UpdateRuleActiveOff   = "Rule deactivated."
	ReorderRulesSuccess   = "Rules reordered successfully."
	DeleteRuleSuccess     = "Rule deleted successfully."
	ErrRuleActiveRequired = "The active field is required."
	ErrRuleActiveInvalid  = "The active field must be true or false."
)

type RuleService interface {
	GetRules(context.Context, string) ([]model.Rule, error)
	PostRule(context.Context, model.Rule) (model.Rule, error)
	UpdateRule(context.Context, model.Rule) (model.Rule, error)
	UpdateRuleActive(context.Context, string, string, bool) error
	ReorderRules(context.Context, string, []string) error
	DeleteRule(context.Context, string, string) error
}

// @Summary Get rules
// @Description Get all of the authenticated user's rules, in the order they're applied
// @Tags rule
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} model.Rule
// @Failure 400 {object} ErrorRes
// @Router /rules [get]
func (h *Handler) GetRules(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)

	rules, err := h.Service.GetRules(c.Context(), userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(rules)
}

// @Summary Create rule
// @Description Create a rule, applied after the user's existing rules
// @Tags rule
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body RuleReq true "Rule request"
// @Success 201 {object} map[string]any "message, rule"
// @Failure 400 {object} ErrorRes
// @Router /rule [post]
func (h *Handler) PostRule(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	rule, ok := h.parseRule(c, userID)
	if !ok {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	created, err := h.Service.PostRule(c.Context(), rule)
	if err != nil {
		return ruleError(c, err)
	}

	return c.Status(201).JSON(fiber.Map{
		"message": PostRuleSuccess,
		"rule":    created,
	})
}

// @Summary Update rule
// @Description Replace a rule's name, conditions, actions and email types. Its order, active flag and Applied count are kept.
// @Tags rule
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Rule ID"
// @Param body body RuleReq true "Rule request"
// @Success 200 {object} map[string]any "message, rule"
// @Failure 400 {object} ErrorRes
// @Failure 404 {object} ErrorRes
// @Router /rule/{id} [put]
func (h *Handler) UpdateRule(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	rule, ok := h.parseRule(c, userID)
	if !ok {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}
	rule.ID = c.Params("id")

	updated, err := h.Service.UpdateRule(c.Context(), rule)
	if err != nil {
		return ruleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": UpdateRuleSuccess,
		"rule":    updated,
	})
}

// @Summary Update rule active
// @Description Turn a rule on or off
// @Tags rule
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Rule ID"
// @Param body body RuleActiveReq true "Rule active request"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Failure 404 {object} ErrorRes
// @Router /rule/{id}/active [put]
func (h *Handler) UpdateRuleActive(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	req := RuleActiveReq{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrRuleActiveInvalid,
		})
	}
	if req.Active == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrRuleActiveRequired,
		})
	}

	err := h.Service.UpdateRuleActive(c.Context(), c.Params("id"), userID, *req.Active)
	if err != nil {
		return ruleError(c, err)
	}

	message := UpdateRuleActiveOff
	if *req.Active {
		message = UpdateRuleActiveOn
	}

	return c.JSON(fiber.Map{
		"message": message,
	})
}

// @Summary Reorder rules
// @Description Set the order rules are applied in. ids must be exactly the user's rule IDs.
// @Tags rule
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body RuleOrderReq true "Rule order request"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /rules/order [put]
func (h *Handler) ReorderRules(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	req := RuleOrderReq{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}
	if err := h.Validator.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

	err := h.Service.ReorderRules(c.Context(), userID, req.IDs)
	if err != nil {
		return ruleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": ReorderRulesSuccess,
	})
}

// @Summary Delete rule
// @Description Delete a rule
// @Tags rule
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Rule ID"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Failure 404 {object} ErrorRes
// @Router /rule/{id} [delete]
func (h *Handler) DeleteRule(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)

	err := h.Service.DeleteRule(c.Context(), c.Params("id"), userID)
	if err != nil {
		return ruleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": DeleteRuleSuccess,
	})
}

// parseRule reads a rule definition from the request body. ok is false if it's malformed.
func (h *Handler) parseRule(c *fiber.Ctx, userID string) (model.Rule, bool) {
	req := RuleReq{}
	if err := c.BodyParser(&req); err != nil {
		return model.Rule{}, false
	}
	if err := h.Validator.Struct(req); err != nil {
		return model.Rule{}, false
	}

	rule := model.Rule{
		UserID:     userID,
		Name:       req.Name,
		Operator:   model.RuleOperator(req.Operator),
		Conditions: make([]model.RuleCondition, 0, len(req.Conditions)),
		Actions:    make([]model.RuleAction, 0, len(req.Actions)),
		Forwards:   req.Forwards,
		Replies:    req.Replies,
		Sends:      req.Sends,
	}
	for _, cond := range req.Conditions {
		rule.Conditions = append(rule.Conditions, model.RuleCondition{
			Type:   model.RuleConditionType(cond.Type),
			Match:  model.RuleMatch(cond.Match),
			Values: cond.Values,
		})
	}
	for _, action := range req.Actions {
		rule.Actions = append(rule.Actions, model.RuleAction{
			Type:  model.RuleActionType(action.Type),
			Value: action.Value,
		})
	}

	return rule, true
}

// ruleError responds with a service error: 404 for a rule that doesn't exist or isn't the
// user's, 400 for anything else.
func ruleError(c *fiber.Ctx, err error) error {
	status := 400
	if errors.Is(err, model.ErrRuleNotFound) {
		status = 404
	}

	return c.Status(status).JSON(fiber.Map{
		"error": err.Error(),
	})
}
