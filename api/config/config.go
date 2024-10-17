package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type APIConfig struct {
	FQDN              string
	Name              string
	Port              string
	ApiAllowOrigin    string
	ApiAllowIp        string
	CookieSecret      string
	TokenSecret       string
	TokenExpiration   time.Duration
	PSK               string
	PSKAllowOrigin    string
	Domains           string
	LogFile           string
	BasicAuthUser     string
	BasicAuthPassword string
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
	Host       string
	Port       string
	User       string
	Password   string
	Sender     string
	SenderName string
}

type ServiceConfig struct {
	OTPExpiration     time.Duration
	SubscriptionType  string
	MaxCredentials    int
	MaxRecipients     int
	MaxDailyAliases   int
	MaxDailySendReply int
}

type Config struct {
	API        APIConfig
	DB         DBConfig
	Redis      RedisConfig
	SMTPClient SMTPClientConfig
	Service    ServiceConfig
}

func New() (Config, error) {
	err := godotenv.Load("./.env")
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

	maxRecipients, err := strconv.Atoi(os.Getenv("MAX_RECIPIENTS"))
	if err != nil {
		return Config{}, err
	}

	maxCredentials, err := strconv.Atoi(os.Getenv("MAX_CREDENTIALS"))
	if err != nil {
		return Config{}, err
	}

	maxDailyAliases, err := strconv.Atoi(os.Getenv("MAX_DAILY_ALIASES"))
	if err != nil {
		return Config{}, err
	}

	maxDailySendReply, err := strconv.Atoi(os.Getenv("MAX_DAILY_SEND_REPLY"))
	if err != nil {
		return Config{}, err
	}

	return Config{
		API: APIConfig{
			FQDN:              os.Getenv("FQDN"),
			Name:              os.Getenv("API_NAME"),
			Port:              os.Getenv("API_PORT"),
			ApiAllowOrigin:    os.Getenv("API_ALLOW_ORIGIN"),
			ApiAllowIp:        os.Getenv("API_ALLOW_IP"),
			CookieSecret:      os.Getenv("COOKIE_SECRET"),
			TokenSecret:       os.Getenv("TOKEN_SECRET"),
			TokenExpiration:   tokenExp,
			PSK:               os.Getenv("PSK"),
			PSKAllowOrigin:    os.Getenv("PSK_ALLOW_ORIGIN"),
			Domains:           os.Getenv("DOMAINS"),
			LogFile:           os.Getenv("LOG_FILE"),
			BasicAuthUser:     os.Getenv("BASIC_AUTH_USER"),
			BasicAuthPassword: os.Getenv("BASIC_AUTH_PASSWORD"),
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
			Host:       os.Getenv("SMTP_CLIENT_HOST"),
			Port:       os.Getenv("SMTP_CLIENT_PORT"),
			User:       os.Getenv("SMTP_CLIENT_USER"),
			Password:   os.Getenv("SMTP_CLIENT_PASSWORD"),
			Sender:     os.Getenv("SMTP_CLIENT_SENDER"),
			SenderName: os.Getenv("SMTP_CLIENT_SENDER_NAME"),
		},

		Service: ServiceConfig{
			OTPExpiration:     otpExp,
			SubscriptionType:  os.Getenv("SUBSCRIPTION_TYPE"),
			MaxCredentials:    maxCredentials,
			MaxRecipients:     maxRecipients,
			MaxDailyAliases:   maxDailyAliases,
			MaxDailySendReply: maxDailySendReply,
		},
	}, nil
}
