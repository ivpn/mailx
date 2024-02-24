package service

import (
	"context"
	"time"
)

type Store interface {
	RecipienteStore
	AliasStore
	UserStore
}

type Cache interface {
	Set(context.Context, string, interface{}, time.Duration) error
	Get(context.Context, string) (string, error)
}

type Service struct {
	Store Store
	Cache Cache
}

func New(store Store, cache Cache) *Service {
	return &Service{
		Store: store,
		Cache: cache,
	}
}
