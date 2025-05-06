package http

import (
	"bytes"
	"encoding/json"
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

func (h Http) SignupWebhookNet(subID string) error {
	data := map[string]string{
		"uuid": subID,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error calling signup webhook:", err)
		return err
	}

	req, err := http.NewRequest("POST", h.Cfg.SignupWebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error calling signup webhook:", err)
		return err
	}

	req.Header.Set("Authorization", "Bearer "+h.Cfg.SignupWebhookPSK)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error calling signup webhook:", err)
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (h Http) SignupWebhook(subID string) error {
	req := fiber.Post(h.Cfg.SignupWebhookURL)
	req.Set("Content-Type", "application/json")
	req.Set("Authorization", "Bearer "+h.Cfg.SignupWebhookPSK)
	req.Body([]byte(`{"uuid": "` + subID + `"}`))

	status, body, err := req.Bytes()
	if err != nil {
		log.Println("Error calling signup webhook:", err)
		return errors.New("Error calling signup webhook")
	}

	log.Println("Signup webhook request:", req)

	if status != http.StatusOK {
		log.Println("Error calling signup webhook, status:", status)
		return errors.New("Error calling signup webhook")
	}

	log.Println("Signup webhook response:", string(body))

	return nil
}
