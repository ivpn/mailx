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
	"github.com/golang-jwt/jwt"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/model"
)

var (
	ErrNoToken  = fmt.Errorf("no token")
	ErrNoClaims = fmt.Errorf("no claims")
	ErrNoExp    = fmt.Errorf("no exp")
)

const (
	AUTH_COOKIE  = "auth"
	AUTHN_COOKIE = "authn"
	USER_ID      = "user_id"
)

type Cache interface {
	Get(context.Context, string) (string, error)
}

type Service interface {
	GetSession(context.Context, string) (model.Session, bool, error)
	GetUser(context.Context, string) (model.User, error)
}

func New(cfg config.APIConfig, cache Cache) fiber.Handler {

	return func(c *fiber.Ctx) error {
		jwtString := GetToken(c)
		if jwtString == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		jwtSignature := GetTokenSignature(jwtString)
		jwtInvalid, _ := cache.Get(c.Context(), "logout_"+jwtSignature)
		if jwtInvalid != "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(cfg.TokenSecret), nil
		})
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		if claims[USER_ID] == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Locals(USER_ID, claims[USER_ID])

		return c.Next()
	}
}

func Webauthn(s Service) fiber.Handler {

	return func(c *fiber.Ctx) error {
		token := c.Cookies(AUTHN_COOKIE)
		if token == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		session, ok, err := s.GetSession(c.Context(), token)
		if err != nil || !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		user, err := s.GetUser(c.Context(), session.UserID)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Locals(USER_ID, user.ID)

		return c.Next()
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

func NewMailserverCORS(cfg config.APIConfig) fiber.Handler {

	return func(c *fiber.Ctx) error {
		if c.IP() != cfg.ApiAllowIp {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.Next()
	}
}

func GetUserID(c *fiber.Ctx) string {
	return c.Locals(USER_ID).(string)
}

func GetToken(c *fiber.Ctx) string {
	var tokenString string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies(AUTH_COOKIE) != "" {
		tokenString = c.Cookies(AUTH_COOKIE)
	}

	return tokenString
}

func GetTokenExp(cfg config.APIConfig, c *fiber.Ctx) (time.Duration, error) {
	jwtString := GetToken(c)
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(cfg.TokenSecret), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrNoClaims
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return 0, ErrNoExp
	}

	return time.Until(time.Unix(int64(exp), 0)), nil
}

func GetTokenSignature(jwt string) string {
	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		return ""
	}

	return parts[2]
}

func NewCookie(token string, cfg config.APIConfig) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     AUTH_COOKIE,
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(time.Duration(cfg.TokenExpiration)),
	}
}

func NewCookieAuthn(token string, path string, cfg config.APIConfig) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     AUTHN_COOKIE,
		Value:    token,
		Path:     path,
		HTTPOnly: true,
		Secure:   true,
		MaxAge:   int(cfg.TokenExpiration),
		Expires:  time.Now().Add(time.Duration(cfg.TokenExpiration)),
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
