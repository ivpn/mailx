package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

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

func (h Http) SignupWebhook(subID string) error {
	req := fiber.Post(h.Cfg.SignupWebhookURL)
	req.Set("Content-Type", "application/json")
	req.Set("Accept", "application/json")
	req.Set("Authorization", "Bearer "+h.Cfg.SignupWebhookPSK)
	req.Body([]byte(`{"uuid": "` + subID + `"}`))

	status, _, err := req.Bytes()
	if err != nil {
		log.Printf("Error calling signup webhook: %v", err)
		return errors.New("error calling signup webhook")
	}

	if status != http.StatusOK {
		log.Printf("Error calling signup webhook, status: %d", status)
		return errors.New("error response from signup webhook")
	}

	return nil
}

func (h Http) GetPreauth(ID string) (model.Preauth, error) {
	req := fiber.Get(h.Cfg.PreauthURL + "/" + ID)
	req.Set("Content-Type", "application/json")
	req.Set("Accept", "application/json")
	req.Set("Authorization", "Bearer "+h.Cfg.PreauthPSK)

	var preauth model.Preauth
	status, body, err := req.Bytes()
	if err != nil {
		log.Printf("Error calling preauth service: %v", err)
		return model.Preauth{}, errors.New("error calling preauth service")
	}

	if status != http.StatusOK {
		log.Printf("Error calling preauth service, status: %d", status)
		return model.Preauth{}, errors.New("error response from preauth service")
	}

	if err := json.Unmarshal(body, &preauth); err != nil {
		log.Printf("Error parsing preauth response: %v", err)
		return model.Preauth{}, errors.New("error parsing preauth response")
	}

	return preauth, nil
}
