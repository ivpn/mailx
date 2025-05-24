package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type APIConfig struct {
	FQDN              string
	Name              string
	Port              string
	ApiAllowOrigin    string
	TokenSecret       string
	TokenExpiration   time.Duration
	PSK               string
	PSKAllowOrigin    string
	Domains           string
	LogFile           string
	BasicAuthUser     string
	BasicAuthPassword string
	SignupWebhookURL  string
	SignupWebhookPSK  string
}

type DBConfig struct {
	Hosts    []string
	Port     string
	Name     string
	User     string
	Password string
}

type RedisConfig struct {
	Addr                  string
	Addrs                 []string
	MasterName            string
	Username              string
	Password              string
	FailoverUsername      string
	FailoverPassword      string
	TLSEnabled            bool
	CertFile              string
	KeyFile               string
	CACertFile            string
	TLSInsecureSkipVerify bool // Optional: Only for testing, use false in production
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
	OTPExpiration          time.Duration
	SubscriptionType       string
	MaxCredentials         int
	MaxRecipients          int
	MaxDailyAliases        int
	MaxDailySendReply      int
	MaxSessions            int
	ForwardGracePeriodDays int
	AccountGracePeriodDays int
	IdLimiterMax           int
	IdLimiterExpiration    time.Duration
}

type Config struct {
	API        APIConfig
	DB         DBConfig
	Redis      RedisConfig
	SMTPClient SMTPClientConfig
	Service    ServiceConfig
}

func New() (Config, error) {
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

	idLimiterMax, err := strconv.Atoi(os.Getenv("ID_LIMITER_MAX"))
	if err != nil {
		return Config{}, err
	}
	idLimiterExpirationStr := os.Getenv("ID_LIMITER_EXPIRATION")
	idLimiterExpiration, err := time.ParseDuration(idLimiterExpirationStr)
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

	maxSessions, err := strconv.Atoi(os.Getenv("MAX_SESSIONS"))
	if err != nil {
		return Config{}, err
	}

	forwardGracePeriodDays, err := strconv.Atoi(os.Getenv("FORWARD_GRACE_PERIOD_DAYS"))
	if err != nil {
		return Config{}, err
	}

	accountGracePeriodDays, err := strconv.Atoi(os.Getenv("ACCOUNT_GRACE_PERIOD_DAYS"))
	if err != nil {
		return Config{}, err
	}

	dbHosts := strings.Split(os.Getenv("DB_HOSTS"), ",")
	redisAddrs := strings.Split(os.Getenv("REDIS_ADDRESSES"), ",")

	return Config{
		API: APIConfig{
			FQDN:              os.Getenv("FQDN"),
			Name:              os.Getenv("API_NAME"),
			Port:              os.Getenv("API_PORT"),
			ApiAllowOrigin:    os.Getenv("API_ALLOW_ORIGIN"),
			TokenSecret:       os.Getenv("TOKEN_SECRET"),
			TokenExpiration:   tokenExp,
			PSK:               os.Getenv("PSK"),
			PSKAllowOrigin:    os.Getenv("PSK_ALLOW_ORIGIN"),
			Domains:           os.Getenv("DOMAINS"),
			LogFile:           os.Getenv("LOG_FILE"),
			BasicAuthUser:     os.Getenv("BASIC_AUTH_USER"),
			BasicAuthPassword: os.Getenv("BASIC_AUTH_PASSWORD"),
			SignupWebhookURL:  os.Getenv("SIGNUP_WEBHOOK_URL"),
			SignupWebhookPSK:  os.Getenv("SIGNUP_WEBHOOK_PSK"),
		},
		DB: DBConfig{
			Hosts:    dbHosts,
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
		Redis: RedisConfig{
			Addr:                  os.Getenv("REDIS_ADDR"),
			Addrs:                 redisAddrs,
			MasterName:            os.Getenv("REDIS_MASTER_NAME"),
			Username:              os.Getenv("REDIS_USERNAME"),
			Password:              os.Getenv("REDIS_PASSWORD"),
			FailoverUsername:      os.Getenv("REDIS_FAILOVER_USERNAME"),
			FailoverPassword:      os.Getenv("REDIS_FAILOVER_PASSWORD"),
			TLSEnabled:            os.Getenv("REDIS_TLS_ENABLED") == "true",
			CertFile:              os.Getenv("REDIS_CERT_FILE"),
			KeyFile:               os.Getenv("REDIS_KEY_FILE"),
			CACertFile:            os.Getenv("REDIS_CA_CERT_FILE"),
			TLSInsecureSkipVerify: os.Getenv("REDIS_TLS_INSECURE_SKIP_VERIFY") == "true",
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
			OTPExpiration:          otpExp,
			SubscriptionType:       os.Getenv("SUBSCRIPTION_TYPE"),
			MaxCredentials:         maxCredentials,
			MaxRecipients:          maxRecipients,
			MaxDailyAliases:        maxDailyAliases,
			MaxDailySendReply:      maxDailySendReply,
			MaxSessions:            maxSessions,
			ForwardGracePeriodDays: forwardGracePeriodDays,
			AccountGracePeriodDays: accountGracePeriodDays,
			IdLimiterMax:           idLimiterMax,
			IdLimiterExpiration:    idLimiterExpiration,
		},
	}, nil
}
