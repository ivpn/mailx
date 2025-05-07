package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/config"
)

type Http struct {
	Cfg config.APIConfig
}

func New(cfg config.APIConfig) *Http {
	return &Http{
		Cfg: cfg,
	}
}

func (h Http) SignupWebhook(subID string) error {
	req := fiber.Post(h.Cfg.SignupWebhookURL)
	req.Set("Content-Type", "application/json")
	req.Set("Accept", "application/json")
	req.Set("Authorization", "Bearer "+h.Cfg.SignupWebhookPSK)
	req.Body([]byte(`{"uuid": "` + subID + `"}`))

	log.Println("SignupWebhookURL:", h.Cfg.SignupWebhookURL)
	log.Println("SignupWebhookPSK:", h.Cfg.SignupWebhookPSK)

	status, body, err := req.Bytes()
	if err != nil {
		log.Printf("Error calling signup webhook: %v", err)
		return errors.New("error calling signup webhook")
	}

	if status != http.StatusOK {
		log.Printf("Error from signup webhook, status: %d, body: %s", status, string(body))
		return errors.New("error response from signup webhook")
	}

	log.Printf("Signup webhook successful, response: %s", string(body))
	return nil
}
