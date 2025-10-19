package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/model"
)

var (
	ErrNoToken  = fmt.Errorf("no token")
	ErrNoClaims = fmt.Errorf("no claims")
	ErrNoExp    = fmt.Errorf("no exp")
)

const (
	AUTH_COOKIE       = "auth"
	AUTHN_COOKIE      = "authn"
	AUTHN_TEMP_COOKIE = "authntemp"
	USER_ID           = "user_id"
)

type Cache interface {
	Get(context.Context, string) (string, error)
}

type Service interface {
	GetSession(context.Context, string) (model.Session, bool, error)
	GetUser(context.Context, string) (model.User, error)
}

func New(cfg config.APIConfig, cache Cache, service Service) fiber.Handler {

	return func(c *fiber.Ctx) error {
		if c.Cookies(AUTHN_COOKIE) != "" {
			session, ok, err := service.GetSession(c.Context(), c.Cookies(AUTHN_COOKIE))
			if err == nil && ok {
				user, err := service.GetUser(c.Context(), session.UserID)
				if err == nil {
					c.Locals(USER_ID, user.ID)
					return c.Next()
				}
			}
		}

		return c.SendStatus(fiber.StatusUnauthorized)
	}
}

func NewPSK(cfg config.APIConfig) fiber.Handler {

	return func(c *fiber.Ctx) error {
		if GetToken(c) != cfg.PSK {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.Next()
	}
}

func NewBasicAuth(cfg config.APIConfig) fiber.Handler {

	return basicauth.New(basicauth.Config{
		Users: map[string]string{
			cfg.BasicAuthUser: cfg.BasicAuthPassword,
		},
	})
}

func NewAPICORS(cfg config.APIConfig) fiber.Handler {

	return cors.New(cors.Config{
		AllowOrigins:     cfg.ApiAllowOrigin,
		AllowCredentials: true,
	})
}

func NewPSKCORS(cfg config.APIConfig) fiber.Handler {

	return cors.New(cors.Config{
		AllowOrigins:     cfg.PSKAllowOrigin,
		AllowMethods:     fiber.MethodPut,
		AllowCredentials: true,
	})
}

func GetUserID(c *fiber.Ctx) string {
	return c.Locals(USER_ID).(string)
}

func GetToken(c *fiber.Ctx) string {
	var tokenString string
	authorization := c.Get("Authorization")

	if after, ok := strings.CutPrefix(authorization, "Bearer "); ok {
		tokenString = after
	} else if c.Cookies(AUTH_COOKIE) != "" {
		tokenString = c.Cookies(AUTH_COOKIE)
	}

	return tokenString
}

func GetAuthnToken(c *fiber.Ctx) string {
	return c.Cookies(AUTHN_COOKIE)
}

func NewCookieAuthn(token string, path string, cfg config.APIConfig) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     AUTHN_COOKIE,
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		MaxAge:   int(cfg.TokenExpiration.Seconds()),
		Expires:  time.Now().Add(time.Duration(cfg.TokenExpiration)),
	}
}

func NewCookieTempAuthn(token string, path string, cfg config.APIConfig) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     AUTHN_TEMP_COOKIE,
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		MaxAge:   int(cfg.TokenExpiration.Seconds()),
		Expires:  time.Now().Add(time.Duration(cfg.TokenExpiration)),
	}
}

func NewCookiePASession(id string) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     "pasession",
		Value:    id,
		HTTPOnly: true,
		Secure:   true,
		MaxAge:   900, // 15 minutes
		Expires:  time.Now().Add(15 * time.Minute),
	}
}

func NewWebAuthn(cfg config.APIConfig) *webauthn.WebAuthn {
	var webAuthn *webauthn.WebAuthn
	config := &webauthn.Config{
		RPDisplayName: cfg.Name,                     // Display Name for your site
		RPID:          cfg.FQDN,                     // Generally the FQDN for your site
		RPOrigins:     []string{cfg.ApiAllowOrigin}, // The origin URLs allowed for WebAuthn requests
	}

	webAuthn, err := webauthn.New(config)
	if err != nil {
		fmt.Println(err)
	}

	return webAuthn
}
