package service

import (
	"context"
	"time"

	"ivpn.net/email/api/config"
)

type Store interface {
	RecipienteStore
	AliasStore
	UserStore
	SubscriptionStore
	MessageStore
	SettingsStore
	SessionStore
}

type Cache interface {
	Set(context.Context, string, interface{}, time.Duration) error
	Get(context.Context, string) (string, error)
	Del(context.Context, string) error
}

type Service struct {
	Cfg   config.Config
	Store Store
	Cache Cache
}

func New(cfg config.Config, store Store, cache Cache) *Service {
	return &Service{
		Cfg:   cfg,
		Store: store,
		Cache: cache,
	}
}
