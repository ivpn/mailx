package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerHost     string
	ServerPort     string
	ServerHostname string
}

func New() (Config, error) {
	err := godotenv.Load("./services/inbound/.env")
	cfg := Config{
		ServerHost:     os.Getenv("SERVER_HOST"),
		ServerPort:     os.Getenv("SERVER_PORT"),
		ServerHostname: os.Getenv("SERVER_HOSTNAME"),
	}

	return cfg, err
}
