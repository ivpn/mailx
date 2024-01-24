package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MTAHost     string
	MTAPort     string
	MTAHostname string
}

func New() (Config, error) {
	err := godotenv.Load("./.env")
	cfg := Config{
		MTAHost:     os.Getenv("MTA_HOST"),
		MTAPort:     os.Getenv("MTA_PORT"),
		MTAHostname: os.Getenv("MTA_HOSTNAME"),
	}

	return cfg, err
}
