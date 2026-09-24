package auth

import (
	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/utils/signedurl"
)

// NewSignedURL only lets through requests whose path and query carry a valid, unexpired
// signature. A nil signer (signing not configured) rejects every request.
func NewSignedURL(signer *signedurl.Signer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if signer == nil || !signer.Verify(c.Path(), string(c.Request().URI().QueryString())) {
			return c.SendStatus(fiber.StatusForbidden)
		}

		return c.Next()
	}
}
