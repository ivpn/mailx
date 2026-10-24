package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/model"
)

const maxLoggedResponse = 512

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
	req.Body([]byte(`{"uuid": "` + subID + `", "service": "mail"}`))

	// Log request for debugging
	log.Printf("Signup webhook request: %+v", req)

	status, res, err := req.Bytes()
	if err != nil {
		return fmt.Errorf("error calling signup webhook: %v", err)
	}

	if status != http.StatusOK {
		// Include the response for debugging, capped to keep log lines readable
		if len(res) > maxLoggedResponse {
			res = res[:maxLoggedResponse]
		}
		return fmt.Errorf("error response from signup webhook, status: %d, response: %s", status, string(res))
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
