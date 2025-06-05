package utils

import (
	"context"
	"log"
	"strconv"
	"time"
)

type Cache interface {
	Set(context.Context, string, any, time.Duration) error
	Get(context.Context, string) (string, error)
	Incr(context.Context, string, time.Duration) error
}

type IDLimiter struct {
	ID    string
	Label string
	Max   int
	Exp   time.Duration
	Cache Cache
}

func (l *IDLimiter) Tick() error {
	err := l.Cache.Incr(context.Background(), l.Label+"_"+l.ID, l.Exp)
	if err != nil {
		log.Printf("error setting failed attempts: %s", err.Error())
		return err
	}

	return nil
}

func (l *IDLimiter) IsAllowed() bool {
	failedAttempts, err := l.Cache.Get(context.Background(), l.Label+"_"+l.ID)
	if err != nil {
		failedAttempts = "0"
	}
	failedAttemptsInt, err := strconv.Atoi(failedAttempts)
	if err != nil {
		failedAttemptsInt = 0
	}
	if failedAttemptsInt > l.Max {
		return false
	}

	return true
}
