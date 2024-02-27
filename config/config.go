package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type APIConfig struct {
	Port            string
	CookieSecret    string
	TokenSecret     string
	TokenExpiration time.Duration
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

type RedisConfig struct {
	Addr string
}

type SMTPClientConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Sender   string
}

type ServiceConfig struct {
	OTPExpiration time.Duration
}

type Config struct {
	API        APIConfig
	SMTP       SMTPConfig
	DB         DBConfig
	Redis      RedisConfig
	SMTPClient SMTPClientConfig
	Service    ServiceConfig
}

func New() (Config, error) {
	err := godotenv.Load("./config/.env")
	if err != nil {
		return Config{}, err
	}

	tokenExpStr := os.Getenv("TOKEN_EXPIRATION")
	tokenExp, err := time.ParseDuration(tokenExpStr)
	if err != nil {
		return Config{}, err
	}

	otpExpStr := os.Getenv("OTP_EXPIRATION")
	otpExp, err := time.ParseDuration(otpExpStr)
	if err != nil {
		return Config{}, err
	}

	return Config{
		API: APIConfig{
			Port:            os.Getenv("API_PORT"),
			CookieSecret:    os.Getenv("COOKIE_SECRET"),
			TokenSecret:     os.Getenv("TOKEN_SECRET"),
			TokenExpiration: tokenExp,
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
		Redis: RedisConfig{
			Addr: os.Getenv("REDIS_ADDR"),
		},
		SMTPClient: SMTPClientConfig{
			Host:     os.Getenv("SMTP_CLIENT_HOST"),
			Port:     os.Getenv("SMTP_CLIENT_PORT"),
			User:     os.Getenv("SMTP_CLIENT_USER"),
			Password: os.Getenv("SMTP_CLIENT_PASSWORD"),
			Sender:   os.Getenv("SMTP_CLIENT_SENDER"),
		},
		Service: ServiceConfig{
			OTPExpiration: otpExp,
		},
	}, nil
}
