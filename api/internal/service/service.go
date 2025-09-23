package service

import (
	"context"
	"time"

	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/client/http"
)

type Store interface {
	RecipientsStore
	AliasStore
	UserStore
	SubscriptionStore
	MessageStore
	SettingsStore
	SessionStore
	CredentialStore
	BounceStore
}

type Cache interface {
	Set(context.Context, string, any, time.Duration) error
	Get(context.Context, string) (string, error)
	Del(context.Context, string) error
	Incr(context.Context, string, time.Duration) error
}

type Service struct {
	Cfg   config.Config
	Store Store
	Cache Cache
	Http  http.Http
}

func New(cfg config.Config, store Store, cache Cache) *Service {
	return &Service{
		Cfg:   cfg,
		Store: store,
		Cache: cache,
		Http: http.Http{
			Cfg: cfg.API,
		},
	}
}
