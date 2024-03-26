package utils

import (
	"log"

	"gopkg.in/natefinch/lumberjack.v2"
	"ivpn.net/email-service/config"
)

func NewLogger(cfg config.APIConfig) {
	log.SetOutput(&lumberjack.Logger{
		Filename:   cfg.LogFile,
		MaxBackups: 3,
		MaxAge:     14, //days
	})
}
