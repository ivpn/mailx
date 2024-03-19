package auth

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"ivpn.net/email-service/config"
)

const AuthCookie = "auth"

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

		if claims["user_id"] == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Locals("user_id", claims["user_id"])

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

func GetUserID(c *fiber.Ctx) string {
	return c.Locals("user_id").(string)
}

func getToken(c *fiber.Ctx) string {
	var tokenString string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies(AuthCookie) != "" {
		tokenString = c.Cookies(AuthCookie)
	}

	return tokenString
}
