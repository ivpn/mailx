package http

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/model"
)

type Http struct {
	Cfg config.APIConfig
}

func New(cfg config.APIConfig) *Http {
	return &Http{
		Cfg: cfg,
	}
}

func (http Http) SignupWebhook(sub model.Subscription) error {
	req := fiber.Post(http.Cfg.SignupWebhookURL)
	req.Set("Content-Type", "application/json")
	req.Set("Authorization:", "Bearer "+http.Cfg.SignupWebhookPSK)
	req.Body([]byte(`{"uuid": "` + sub.ID + `"}`))

	_, _, err := req.Bytes()
	if err != nil {
		log.Println("Error calling signup webhook:", err)
		return err[0]
	}

	return nil
}
