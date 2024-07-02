package api

import (
	"bytes"
	"log"

	"github.com/DusanKasan/parsemail"
	"github.com/gofiber/fiber/v2"
)

type ProcessorService interface {
}

func (h *Handler) Health(c *fiber.Ctx) error {
	log.Println("Health check OK")
	return c.SendString("OK")
}

func (h *Handler) HandleEmail(c *fiber.Ctx) error {
	log.Println("Email received", c.Body())

	var reader = bytes.NewReader(c.Body())
	email, err := parsemail.Parse(reader)
	if err != nil {
		log.Println("error parsing email", err)
	}

	log.Println("Email parsed", email.From, email.Subject)

	return c.SendString("OK")
}
