package api

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
)

var (
	RegisterSuccess              = "User is created"
	LoginSuccess                 = "Login successful"
	LogoutSuccess                = "Logout successful"
	DeleteUserSuccess            = "User is deleted"
	OTPSent                      = "New OTP is sent"
	ActivateUserSuccess          = "Email is confirmed"
	ChangePasswordSuccess        = "Password is changed"
	InitiatePasswordResetSuccess = "Password reset initiated"
	ErrInvalidCredentials        = "Invalid credentials"
	ErrInvalidRequest            = "Invalid request"
	ErrLogoutUser                = "Could not logout user"
)

type UserService interface {
	PostUser(context.Context, model.User) error
	SendUserOTP(context.Context, string) error
	ActivateUser(context.Context, string, string) error
	GetUserByCredentials(context.Context, string, string) (model.User, error)
	GetUserByPassword(context.Context, string, string) (model.User, error)
	DeleteUser(context.Context, string) error
	GetUser(context.Context, string) (model.User, error)
	GetUserStats(context.Context, string) (model.UserStats, error)
	LogoutUser(context.Context, string, time.Duration) error
	ChangePassword(context.Context, string, string) error
	InitiatePasswordReset(context.Context, string) error
}

// @Summary Register user
// @Description Register user
// @Tags user
// @Accept json
// @Produce json
// @Param body body UserReq true "User request"
// @Success 201 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /p/register [post]
func (h *Handler) Register(c *fiber.Ctx) error {
	// Parse the request
	req := UserReq{}
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

	// Save the user
	err = h.Service.PostUser(c.Context(), user)
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
// @Router /p/login [post]
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

	// Create auth token
	token, err := utils.CreateAuthToken(h.Cfg, user.ID, user.Email)
	if err != nil {
		log.Printf("error login: %s", err.Error())
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidCredentials,
		})
	}

	// Set the token in encrypted cookie
	c.Cookie(auth.NewCookie(token, h.Cfg))

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
	jwt := auth.GetToken(c)
	jwtSignature := auth.GetTokenSignature(jwt)
	jwtExp, err := auth.GetTokenExp(h.Cfg, c)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": ErrLogoutUser,
		})
	}

	err = h.Service.LogoutUser(c.Context(), jwtSignature, jwtExp)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.ClearCookie(auth.AUTH_COOKIE)

	return c.Status(200).JSON(fiber.Map{
		"message": LogoutSuccess,
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

	// Get the user
	user, err := h.Service.GetUserByPassword(c.Context(), ID, req.Password)
	if err != nil || user.ID != ID {
		log.Printf("error deleting user: %s", err.Error())
		return c.Status(401).JSON(fiber.Map{
			"error": ErrInvalidCredentials,
		})
	}

	// Delete the user
	err = h.Service.DeleteUser(c.Context(), ID)
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
		"message": ChangePasswordSuccess,
	})
}

// @Summary Initiate password reset
// @Description Initiate password reset
// @Tags user
// @Accept json
// @Produce json
// @Param body body InitiatePasswordResetReq true "Initiate password reset request"
// @Success 200 {object} SuccessRes
// @Failure 400 {object} ErrorRes
// @Router /p/initiatepasswordreset [post]
func (h *Handler) InitiatePasswordReset(c *fiber.Ctx) error {
	req := InitiatePasswordResetReq{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Printf("error initiating password reset: %s", err.Error())
		return c.Status(400).JSON(fiber.Map{
			"error": ErrInvalidRequest,
		})
	}

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
