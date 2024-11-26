package utils

import (
	"log"

	"gopkg.in/natefinch/lumberjack.v2"
	"ivpn.net/email/api/config"
)

func NewLogger(cfg config.APIConfig) {
	log.SetOutput(&lumberjack.Logger{
		Filename:   cfg.LogFile,
		MaxBackups: 3,
		MaxAge:     7, //days
	})
}
