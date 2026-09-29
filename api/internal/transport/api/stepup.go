package api

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
)

var StepUpSuccess = "Verification successful."

type StepUpService interface {
	SetStepUp(context.Context, string) error
	ConsumeStepUp(context.Context, string) (bool, error)
}

// @Summary Step-up verification via password
// @Description Verify the current user's password to authorize a pending sensitive action
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body StepUpPasswordReq true "Step-up password request"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /user/stepup/password [post]
func (h *Handler) StepUpPassword(c *fiber.Ctx) error {
	req := StepUpPasswordReq{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Printf("error step-up password: %s", err.Error())
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

	ID := auth.GetUserID(c)
	_, err = h.Service.GetUserByPassword(c.Context(), ID, req.Password)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = h.Service.SetStepUp(c.Context(), auth.GetAuthnCookie(c))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": StepUpSuccess,
	})
}

// @Summary Begin step-up verification via passkey
// @Description Begin a WebAuthn assertion ceremony to authorize a pending sensitive action
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /user/stepup/passkey/begin [post]
func (h *Handler) StepUpPasskeyBegin(c *fiber.Ctx) error {
	ID := auth.GetUserID(c)

	user, err := h.Service.GetUser(c.Context(), ID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	options, sessionData, err := h.WebAuthn.BeginLogin(user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	exp := time.Now().Add(auth.WebAuthnCeremonyExpiration)
	token, err := model.GenSessionToken()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrSaveSession,
		})
	}
	sessionData.Expires = exp
	err = h.Service.SaveSession(c.Context(), *sessionData, token, user.ID, exp, false)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrSaveSession,
		})
	}

	c.Cookie(auth.NewCookieTempAuthn(token, c.Path(), h.Cfg))

	return c.Status(200).JSON(options)
}

// @Summary Finish step-up verification via passkey
// @Description Finish a WebAuthn assertion ceremony to authorize a pending sensitive action
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /user/stepup/passkey/finish [post]
func (h *Handler) StepUpPasskeyFinish(c *fiber.Ctx) error {
	token := c.Cookies(auth.AUTHN_TEMP_COOKIE)

	session, ok, err := h.Service.GetSession(c.Context(), token)
	if err != nil || !ok {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrGetSession,
		})
	}

	// Ensure the ceremony belongs to the already-logged-in caller
	if auth.GetUserID(c) != session.UserID {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrGetSession,
		})
	}

	user, err := h.Service.GetUser(c.Context(), session.UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	r, err := adaptor.ConvertRequest(c, true)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrFinishLogin,
		})
	}

	credential, err := h.WebAuthn.FinishLogin(user, session.SessionData, r)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if credential.Authenticator.CloneWarning {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrFinishLogin,
		})
	}

	err = h.Service.UpdateCredential(c.Context(), *credential, user.ID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = h.Service.DeleteSession(c.Context(), token)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrDeleteSession,
		})
	}
	auth.ClearCookies(c, auth.AUTHN_TEMP_COOKIE)

	err = h.Service.SetStepUp(c.Context(), auth.GetAuthnCookie(c))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": StepUpSuccess,
	})
}
