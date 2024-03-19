package auth

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"ivpn.net/email-service/config"
)

const (
	AUTH_COOKIE = "auth"
	USER_ID     = "user_id"
)

func New(cfg config.APIConfig) fiber.Handler {

	return func(c *fiber.Ctx) error {
		tokenString := getToken(c)
		if tokenString == "" {
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
		if getToken(c) != cfg.PSK {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.Next()
	}
}

func NewPSKCORS(cfg config.APIConfig) fiber.Handler {

	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", cfg.PSKAllowOrigin)
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		c.Set("Access-Control-Allow-Methods", "PUT")

		return c.Next()
	}
}

func GetUserID(c *fiber.Ctx) string {
	return c.Locals(USER_ID).(string)
}

func getToken(c *fiber.Ctx) string {
	var tokenString string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies(AUTH_COOKIE) != "" {
		tokenString = c.Cookies(AUTH_COOKIE)
	}

	return tokenString
}
