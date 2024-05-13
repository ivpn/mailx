package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt"
	"ivpn.net/email/api/config"
)

var (
	ErrNoToken  = fmt.Errorf("no token")
	ErrNoClaims = fmt.Errorf("no claims")
	ErrNoExp    = fmt.Errorf("no exp")
)

const (
	AUTH_COOKIE = "auth"
	USER_ID     = "user_id"
)

type Cache interface {
	Get(context.Context, string) (string, error)
}

func New(cfg config.APIConfig, cache Cache) fiber.Handler {

	return func(c *fiber.Ctx) error {
		tokenString := GetToken(c)
		if tokenString == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		tokenHash, err := argon2id.CreateHash(tokenString, argon2id.DefaultParams)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		loggedOut, _ := cache.Get(c.Context(), "logout_"+tokenHash)
		if loggedOut != "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies(AUTH_COOKIE) != "" {
		tokenString = c.Cookies(AUTH_COOKIE)
	}

	return tokenString
}

func GetTokenExp(c *fiber.Ctx) (time.Duration, error) {
	tokenString := GetToken(c)
	if tokenString == "" {
		return 0, ErrNoToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
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
