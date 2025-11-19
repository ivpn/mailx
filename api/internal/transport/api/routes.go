package api

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/swagger"
	"ivpn.net/email/api/config"
	_ "ivpn.net/email/api/docs"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/middleware/limit"
)

func (h *Handler) SetupRoutes(cfg config.APIConfig) {
	email := h.Server.Group("/v1/email")
	email.Use(auth.NewPSKCORS(cfg))
	email.Use(auth.NewPSK(cfg))
	email.Post("", h.HandleEmail)

	h.Server.Use(auth.NewAPICORS(cfg))
	h.Server.Use(helmet.New())
	h.Server.Use(healthcheck.New())

	h.Server.Post("/v1/register", limiter.New(), h.Register)
	h.Server.Post("/v1/login", limit.New(5, 10*time.Minute), h.Login)
	h.Server.Post("/v1/initiatepasswordreset", limiter.New(), h.InitiatePasswordReset)
	h.Server.Put("/v1/resetpassword", limiter.New(), h.ResetPassword)
	h.Server.Put("/v1/rotatepasession", limiter.New(), h.RotatePASession)

	h.Server.Post("/v1/register/begin", limiter.New(), h.BeginRegistration)
	h.Server.Post("/v1/register/finish", limiter.New(), h.FinishRegistration)
	h.Server.Post("/v1/login/begin", limiter.New(), h.BeginLogin)
	h.Server.Post("/v1/login/finish", limiter.New(), h.FinishLogin)

	session := h.Server.Group("/v1/pasession")
	session.Use(auth.NewPSK(cfg))
	session.Post("/add", h.AddPASession)

	v1 := h.Server.Group("/v1")
	v1.Use(auth.New(cfg, h.Cache, h.Service))

	v1.Post("/register/add", limiter.New(), h.AddPasskey)
	v1.Post("/register/add/finish", limiter.New(), h.FinishAddPasskey)
	v1.Post("/user/sendotp", limit.New(5, 10*time.Minute), h.SendUserOTP)
	v1.Post("/user/activate", limiter.New(), h.Activate)
	v1.Post("/user/logout", h.Logout)
	v1.Post("/user/delete/request", limit.New(5, 10*time.Minute), h.DeleteUserRequest)
	v1.Post("/user/delete", limit.New(5, 10*time.Minute), h.DeleteUser)
	v1.Get("/user", h.GetUser)
	v1.Get("/user/stats", h.GetUserStats)
	v1.Get("/user/credentials", h.GetCredentials)
	v1.Delete("/user/credential/:id", h.DeleteCredential)
	v1.Put("/user/changepassword", limit.New(5, 10*time.Minute), h.ChangePassword)
	v1.Put("/user/changeemail", limit.New(5, 10*time.Minute), h.ChangeEmail)
	v1.Put("/user/totp/enable", limit.New(5, 10*time.Minute), h.TotpEnable)
	v1.Put("/user/totp/enable/confirm", limit.New(5, 10*time.Minute), h.TotpEnableConfirm)
	v1.Put("/user/totp/disable", limit.New(5, 10*time.Minute), h.TotpDisable)

	v1.Get("/sub", h.GetSubscription)
	v1.Put("/sub/update", limiter.New(), h.UpdateSubscription)

	v1.Get("/settings", h.GetSettings)
	v1.Put("/settings", h.UpdateSettings)

	v1.Get("/recipient/:id", h.GetRecipient)
	v1.Get("/recipients", h.GetRecipients)
	v1.Post("/recipient", limit.New(5, 10*time.Minute), h.PostRecipient)
	v1.Put("/recipient", h.UpdateRecipient)
	v1.Post("/recipient/sendotp/:id", limit.New(5, 10*time.Minute), h.SendRecipientOTP)
	v1.Post("/recipient/activate/:id", limit.New(5, 10*time.Minute), h.ActivateRecipient)
	v1.Put("/recipient/delete/:id", h.DeleteRecipient)

	v1.Get("/alias/:id", h.GetAlias)
	v1.Get("/aliases", h.GetAliases)
	v1.Get("/aliases/export", h.ExportAliases)
	v1.Post("/alias", limiter.New(), h.PostAlias)
	v1.Put("/alias/:id", h.UpdateAlias)
	v1.Delete("/alias/:id", h.DeleteAlias)

	v1.Get("/logs", h.GetLogs)
	v1.Get("/log/file/:id", h.GetLogFile)

	docs := h.Server.Group("/docs")
	docs.Use(auth.NewBasicAuth(cfg))
	docs.Get("/*", swagger.HandlerDefault)
}
