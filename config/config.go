package config

import (
	"os"

	"github.com/joho/godotenv"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Hostname string
}

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type Config struct {
	SMTP SMTPConfig
	DB   DBConfig
}

func New() (Config, error) {
	err := godotenv.Load("./config/.env")
	if err != nil {
		return Config{}, err
	}

	return Config{
		SMTP: SMTPConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     os.Getenv("SMTP_PORT"),
			Hostname: os.Getenv("SMTP_HOSTNAME"),
		},
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
	}, nil
}
