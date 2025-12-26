package api

import (
	"context"
	"log"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
)

var (
	RegisterSuccess              = "Account created successfully."
	LoginSuccess                 = "Logged in successfully."
	LogoutSuccess                = "Logged out successfully."
	DeleteUserSuccess            = "Account deleted successfully."
	OTPSent                      = "A new OTP has been sent to your email."
	ActivateUserSuccess          = "Email verified successfully."
	InitiatePasswordResetSuccess = "Password reset link sent to your email."
	ResetPasswordSuccess         = "Your password has been changed successfully."
	ChangeEmailSuccess           = "Your email has been updated. A 6-digit OTP code has been sent to your new email for verification."
	ErrInvalidCredentials        = "The email or password you entered is incorrect."
	ErrInvalidRequest            = "The request is invalid. Please check your input and try again."
	ErrLogoutUser                = "Could not log out. Please try again."
	DisableTotpSuccess           = "Two-factor authentication has been disabled."
	TotpRequired                 = "Two-factor authentication is required to continue."
	ErrInvalidTotpCode           = "The 2FA code you entered is invalid."
	ErrGetUser                   = "We couldn’t retrieve your user details."
	ErrTooManySessions           = "You have too many active sessions. Please log out from other devices or try again later."
)

type UserService interface {
	SendUserOTP(context.Context, string) error
	ActivateUser(context.Context, string, string) error
	GetUserByCredentials(context.Context, string, string) (model.User, error)
	GetUserByPassword(context.Context, string, string) (model.User, error)
	GetUserByEmail(context.Context, string) (model.User, error)
	GetUnfinishedSignupOrPostUser(context.Context, model.User, string, string) (model.User, error)
	SaveUser(context.Context, model.User) error
	DeleteUserRequest(context.Context, string) (string, error)
	DeleteUser(context.Context, string, string) error
	GetUser(context.Context, string) (model.User, error)
	GetUserStats(context.Context, string) (model.UserStats, error)
	LogoutUser(context.Context, string) error
	ChangePassword(context.Context, string, string) error
	ChangeEmail(context.Context, string, string) error
	InitiatePasswordReset(context.Context, string) error
	ResetPassword(context.Context, string, string) error
	TotpEnable(context.Context, string) (model.TOTPNew, error)
	TotpEnableConfirm(context.Context, string, string) (model.TOTPBackup, error)
	TotpDisable(context.Context, string, string) error
	VerifyTotp(context.Context, string, string) (bool, error)
	CheckSessionCount(context.Context, string) (bool, error)
}

// @Summary Register user
// @Description Register user
// @Tags user
// @Accept json
// @Produce json
// @Param body body SignupUserReq true "User request"
// @Success 201 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /register [post]
func (h *Handler) Register(c *fiber.Ctx) error {
	// Get session ID from cookie
	sessionId := c.Cookies(auth.PA_SESSION_COOKIE)

	// Parse the request
	req := SignupUserReq{}
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
		Email:         req.Email,
		PasswordPlain: &req.Password,
		IsActive:      false,
	}

	// Get unfinished signup user or create new user
	user, err = h.Service.GetUnfinishedSignupOrPostUser(c.Context(), user, req.SubID, sessionId)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get user
	user, err = h.Service.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Set password
	err = h.Service.ChangePassword(c.Context(), user.ID, req.Password)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Send OTP
	err = h.Service.SendUserOTP(c.Context(), user.ID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": RegisterSuccess,
	})
}

// @Summary Send user OTP
// @Description Send user OTP
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /user/sendotp [post]
func (h *Handler) SendUserOTP(c *fiber.Ctx) error {
	ID := auth.GetUserID(c)

	err := h.Service.SendUserOTP(c.Context(), ID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": OTPSent,
	})
}

// @Summary Activate user
// @Description Activate user
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body ActivateReq true "Activate request"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /user/activate [post]
func (h *Handler) Activate(c *fiber.Ctx) error {
	ID := auth.GetUserID(c)

	req := ActivateReq{}
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

	err = h.Service.ActivateUser(c.Context(), ID, req.OTP)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": ActivateUserSuccess,
	})
}

// @Summary Login user
// @Description Login user
// @Tags user
// @Accept json
// @Produce json
// @Param body body UserReq true "User request"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
	// Parse the request
	req := UserReq{}
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

	// Get the user
	user, err := h.Service.GetUserByCredentials(c.Context(), req.Email, req.Password)
	if err != nil {
		log.Printf("error login: %s", err.Error())
		return c.Status(401).JSON(fiber.Map{
			"error": ErrInvalidCredentials,
		})
	}

	// Check max sessions limit
	ok, err := h.Service.CheckSessionCount(c.Context(), user.ID)
	if !ok || err != nil {
		log.Printf("error login: %s", err.Error())
		return c.Status(400).JSON(fiber.Map{
			"error": ErrTooManySessions,
		})
	}

	// Check if TOTP is required
	if user.IsTotpEnabled() && req.OTP == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": TotpRequired,
			"code":    70001,
		})
	}

	// Verify TOTP
	if user.IsTotpEnabled() && req.OTP != "" {
		isValid, err := h.Service.VerifyTotp(c.Context(), user.ID, req.OTP)
		if err != nil {
			log.Printf("error login: %s", err.Error())
		}
		if err != nil || !isValid {
			return c.Status(400).JSON(fiber.Map{
				"error": ErrInvalidTotpCode,
			})
		}
	}

	// Save the session
	exp := time.Now().Add(h.Cfg.TokenExpiration)
	sessionData := webauthn.SessionData{
		UserID:  user.WebAuthnID(),
		Expires: exp,
	}
	token, err := model.GenSessionToken()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrSaveSession,
		})
	}
	err = h.Service.SaveSession(c.Context(), sessionData, token, user.ID, exp)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrSaveSession,
		})
	}

	// Set token in cookie
	c.Cookie(auth.NewCookieAuthn(token, "/", h.Cfg))

	return c.Status(200).JSON(fiber.Map{
		"message": LoginSuccess,
	})
}

// @Summary Logout user
// @Description Logout user
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /user/logout [post]
func (h *Handler) Logout(c *fiber.Ctx) error {
	c.ClearCookie(auth.AUTHN_COOKIE)
	c.ClearCookie(auth.AUTHN_TEMP_COOKIE)

	authnToken := auth.GetAuthnToken(c)

	err := h.Service.LogoutUser(c.Context(), authnToken)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": LogoutSuccess,
	})
}

// @Summary Delete user request
// @Description Delete user request
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /user/delete/request [post]
func (h *Handler) DeleteUserRequest(c *fiber.Ctx) error {
	ID := auth.GetUserID(c)

	otp, err := h.Service.DeleteUserRequest(c.Context(), ID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"otp": otp,
	})
}

// @Summary Delete user
// @Description Delete user
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body DeleteUserReq true "Delete user request"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /user/delete [post]
func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	ID := auth.GetUserID(c)

	// Parse the request
	req := DeleteUserReq{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Printf("error deleting user: %s", err.Error())
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

	// Delete the user
	err = h.Service.DeleteUser(c.Context(), ID, req.OTP)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": DeleteUserSuccess,
	})
}

// @Summary Get user
// @Description Get user
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} model.User
// @Failure 400 {object} ErrorRes
// @Router /user [get]
func (h *Handler) GetUser(c *fiber.Ctx) error {
	ID := auth.GetUserID(c)

	user, err := h.Service.GetUser(c.Context(), ID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(user)
}

// @Summary Get user stats
// @Description Get user stats
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} model.UserStats
// @Failure 400 {object} ErrorRes
// @Router /user/stats [get]
func (h *Handler) GetUserStats(c *fiber.Ctx) error {
	ID := auth.GetUserID(c)

	stats, err := h.Service.GetUserStats(c.Context(), ID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(stats)
}

// @Summary Change password
// @Description Change password
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body ChangePasswordReq true "Change password request"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /user/changepassword [put]
func (h *Handler) ChangePassword(c *fiber.Ctx) error {
	ID := auth.GetUserID(c)

	// Parse the request
	req := ChangePasswordReq{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Printf("error changing password: %s", err.Error())
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

	// Change the password
	err = h.Service.ChangePassword(c.Context(), ID, req.Password)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": ResetPasswordSuccess,
	})
}

// @Summary Change email
// @Description Change email
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body EmailReq true "Change email request"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /user/changeemail [put]
func (h *Handler) ChangeEmail(c *fiber.Ctx) error {
	ID := auth.GetUserID(c)

	// Parse the request
	req := EmailReq{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Printf("error changing email: %s", err.Error())
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

	// Change the email
	err = h.Service.ChangeEmail(c.Context(), ID, req.Email)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Send OTP
	err = h.Service.SendUserOTP(c.Context(), ID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": ChangeEmailSuccess,
	})
}

// @Summary Initiate password reset
// @Description Initiate password reset
// @Tags user
// @Accept json
// @Produce json
// @Param body body EmailReq true "Initiate password reset request"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /initiatepasswordreset [post]
func (h *Handler) InitiatePasswordReset(c *fiber.Ctx) error {
	// Parse the request
	req := EmailReq{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Printf("error initiating password reset: %s", err.Error())
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

	// Initiate password reset
	err = h.Service.InitiatePasswordReset(c.Context(), req.Email)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": InitiatePasswordResetSuccess,
	})
}

// @Summary Reset password
// @Description Reset password
// @Tags user
// @Accept json
// @Produce json
// @Param body body ResetPasswordReq true "Password reset request"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /resetpassword [put]
func (h *Handler) ResetPassword(c *fiber.Ctx) error {
	// Parse the request
	req := ResetPasswordReq{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Printf("error resetting password: %s", err.Error())
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

	// Reset the password
	err = h.Service.ResetPassword(c.Context(), req.OTP, req.Password)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": ResetPasswordSuccess,
	})
}

// @Summary Enable TOTP
// @Description Enable TOTP
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} model.TOTPNew
// @Failure 400 {object} ErrorRes
// @Router /user/totp/enable [put]
func (h *Handler) TotpEnable(c *fiber.Ctx) error {
	// Enable TOTP
	ID := auth.GetUserID(c)
	res, err := h.Service.TotpEnable(c.Context(), ID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(res)
}

// @Summary Enable TOTP confirm
// @Description Enable TOTP confirm
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body TotpReq true "TOTP confirm request"
// @Success 200 {object} model.TOTPBackup
// @Failure 400 {object} ErrorRes
// @Router /user/totp/enable/confirm [put]
func (h *Handler) TotpEnableConfirm(c *fiber.Ctx) error {
	// Parse the request
	req := TotpReq{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Printf("error enabling totp: %s", err.Error())
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

	// Confirm the TOTP
	ID := auth.GetUserID(c)
	res, err := h.Service.TotpEnableConfirm(c.Context(), ID, req.OTP)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(res)
}

// @Summary Disable TOTP
// @Description Disable TOTP
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body TotpReq true "TOTP confirm request"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /user/totp/disable [put]
func (h *Handler) TotpDisable(c *fiber.Ctx) error {
	// Parse the request
	req := TotpReq{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Printf("error enabling totp: %s", err.Error())
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

	// Disable the TOTP
	ID := auth.GetUserID(c)
	err = h.Service.TotpDisable(c.Context(), ID, req.OTP)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": DisableTotpSuccess,
	})
}
