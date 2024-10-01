package api

import (
	"context"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
)

var (
	BeginRegistrationSuccess  = "Registration started"
	FinishRegistrationSuccess = "Registration finished"
	BeginLoginSuccess         = "Login started"
	FinishLoginSuccess        = "Login finished"
	ErrBeginRegistration      = "could not begin registration"
	ErrFinishRegistration     = "could not finish registration"
	ErrBeginLogin             = "could not begin login"
	ErrFinishLogin            = "could not finish login"
	ErrGetSession             = "could not get session"
	ErrSaveSession            = "could not save session"
	ErrDeleteSession          = "could not delete session"
	ErrDeleteCredential       = "could not delete credential"
	DeleteCredentialSuccess   = "Credential deleted"
)

type SessionService interface {
	GetSession(context.Context, string) (model.Session, bool, error)
	SaveSession(context.Context, webauthn.SessionData, string, string) error
	DeleteSession(context.Context, string) error
}

type CredentialService interface {
	GetCredentials(context.Context, string) ([]model.Credential, error)
	SaveCredential(context.Context, webauthn.Credential, string) error
	UpdateCredential(context.Context, webauthn.Credential, string) error
	DeleteCredential(context.Context, webauthn.Credential, string) error
	DeleteCredentialByID(context.Context, string, string) error
}

// @Summary Begin registration
// @Description Begin registration process
// @Tags webauthn
// @Accept json
// @Produce json
// @Param email body EmailReq true "Email"
// @Success 201 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /register/begin [post]
func (h *Handler) BeginRegistration(c *fiber.Ctx) error {
	// Parse the request
	req := EmailReq{}
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

	// Create new user
	user := model.User{
		Email:    req.Email,
		IsActive: false,
	}

	// Get or create user
	user, err = h.Service.GetOrPostUser(c.Context(), user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Begin registration
	options, sessionData, err := h.WebAuthn.BeginRegistration(user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Save the session
	token := model.GenSessionToken()
	err = h.Service.SaveSession(c.Context(), *sessionData, token, user.ID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrSaveSession,
		})
	}

	// Set token in cookie
	c.Cookie(auth.NewCookieAuthn(token, c.Path(), h.Cfg))

	return c.Status(201).JSON(options)
}

// @Summary Finish registration
// @Description Finish registration process
// @Tags webauthn
// @Accept json
// @Produce json
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /register/finish [post]
func (h *Handler) FinishRegistration(c *fiber.Ctx) error {
	// Get cookie token
	token := c.Cookies(auth.AUTHN_COOKIE)

	// Get session
	session, ok, err := h.Service.GetSession(c.Context(), token)
	if err != nil || !ok {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrGetSession,
		})
	}

	// Get user
	user, err := h.Service.GetUser(c.Context(), session.UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Finish registration
	r, err := adaptor.ConvertRequest(c, true)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrFinishRegistration,
		})
	}

	credential, err := h.WebAuthn.FinishRegistration(user, session.SessionData, r)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Add credential to user
	err = h.Service.SaveCredential(c.Context(), *credential, user.ID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Delete session
	err = h.Service.DeleteSession(c.Context(), token)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrDeleteSession,
		})
	}

	// Clear cookie
	c.ClearCookie(auth.AUTHN_COOKIE)

	return c.Status(200).JSON(fiber.Map{
		"message": FinishRegistrationSuccess,
	})
}

// @Summary Begin login
// @Description Begin login process
// @Tags webauthn
// @Accept json
// @Produce json
// @Param email body EmailReq true "Email"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /login/begin [post]
func (h *Handler) BeginLogin(c *fiber.Ctx) error {
	// Parse the request
	req := EmailReq{}
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

	// Get user
	user, err := h.Service.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrGetUser,
		})
	}

	// Begin login
	options, sessionData, err := h.WebAuthn.BeginLogin(user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Save the session
	token := model.GenSessionToken()
	err = h.Service.SaveSession(c.Context(), *sessionData, token, user.ID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrSaveSession,
		})
	}

	// Set token in cookie
	c.Cookie(auth.NewCookieAuthn(token, c.Path(), h.Cfg))

	return c.Status(200).JSON(options)
}

// @Summary Finish login
// @Description Finish login process
// @Tags webauthn
// @Accept json
// @Produce json
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /login/finish [post]
func (h *Handler) FinishLogin(c *fiber.Ctx) error {
	// Get cookie token
	token := c.Cookies(auth.AUTHN_COOKIE)

	// Get session
	session, ok, err := h.Service.GetSession(c.Context(), token)
	if err != nil || !ok {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrGetSession,
		})
	}

	// Get user
	user, err := h.Service.GetUser(c.Context(), session.UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Finish login
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

	// Update user credential
	err = h.Service.UpdateCredential(c.Context(), *credential, user.ID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Delete session
	err = h.Service.DeleteSession(c.Context(), token)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrDeleteSession,
		})
	}

	// Clear cookie
	c.ClearCookie(auth.AUTHN_COOKIE)

	// Save the session
	sessionData := webauthn.SessionData{
		UserID:  user.WebAuthnID(),
		Expires: time.Now().Add(h.Cfg.TokenExpiration),
	}
	token = model.GenSessionToken()
	err = h.Service.SaveSession(c.Context(), sessionData, token, user.ID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrSaveSession,
		})
	}

	// Set token in cookie
	c.Cookie(auth.NewCookieAuthn(token, "/", h.Cfg))

	return c.Status(200).JSON(fiber.Map{
		"message": FinishLoginSuccess,
	})
}

// @Summary Get credentials
// @Description Get user credentials
// @Tags webauthn
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} []model.Credential
// @Failure 400 {object} ErrorRes
// @Router /user/credentials [get]
func (h *Handler) GetCredentials(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	credentials, err := h.Service.GetCredentials(c.Context(), userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(credentials)
}

// @Summary Delete credential
// @Description Delete credential by ID
// @Tags webauthn
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Credential ID"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /user/credential/{id} [delete]
func (h *Handler) DeleteCredential(c *fiber.Ctx) error {
	userID := auth.GetUserID(c)
	ID := c.Params("id")
	err := h.Service.DeleteCredentialByID(c.Context(), ID, userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": DeleteCredentialSuccess,
	})
}
