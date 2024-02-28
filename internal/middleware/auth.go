package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"ivpn.net/email-service/config"
)

func NewAuth(cfg config.APIConfig) fiber.Handler {

	return func(c *fiber.Ctx) error {
		tokenString := extractToken(c)

		token, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
			if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
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

		if claims["auth_id"] == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.Next()
	}
}

func extractToken(c *fiber.Ctx) string {
	var tokenString string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("token") != "" {
		tokenString = c.Cookies("token")
	}

	return tokenString
}
