package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SMTPHost     string
	SMTPPort     string
	SMTPHostname string
	DBHost       string
	DBPort       string
	DBName       string
	DBUser       string
	DBPassword   string
}

func New() (Config, error) {
	err := godotenv.Load("./config/.env")
	if err != nil {
		return Config{}, err
	}

	return Config{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPHostname: os.Getenv("SMTP_HOSTNAME"),
		DBHost:       os.Getenv("DB_HOST"),
		DBPort:       os.Getenv("DB_PORT"),
		DBName:       os.Getenv("DB_NAME"),
		DBUser:       os.Getenv("DB_USER"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
	}, nil
}
