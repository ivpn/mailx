package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"ivpn.net/email-service/config"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis() (*Redis, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
	})

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &Redis{
		Client: client,
	}, nil
}

func (r *Redis) Close() error {
	return r.Client.Close()
}

func (r *Redis) Set(key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(context.Background(), key, value, expiration).Err()
}

func (r *Redis) Get(key string) (string, error) {
	return r.Client.Get(context.Background(), key).Result()
}
