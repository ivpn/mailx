package config

import (
	"os"

	"github.com/joho/godotenv"
)

type APIConfig struct {
	Port string
}

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
	API  APIConfig
	SMTP SMTPConfig
	DB   DBConfig
}

func New() (Config, error) {
	err := godotenv.Load("./config/.env")
	if err != nil {
		return Config{}, err
	}

	return Config{
		API: APIConfig{
			Port: os.Getenv("API_PORT"),
		},
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
