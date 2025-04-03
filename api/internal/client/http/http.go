package http

import (
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

func (http Http) PostSubscription(sub model.Subscription) error {
	req := fiber.Post("https://api.example.net/subscription")
	req.Set("Content-Type", "application/json")
	req.Set("Authorization:", "Bearer "+http.Cfg.PSK)
	req.Body([]byte(`{"id": "` + sub.ID + `"}`))

	_, _, err := req.Bytes()
	if err != nil {
		return err[0]
	}

	return nil
}
