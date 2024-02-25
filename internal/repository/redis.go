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

func NewRedis(cfg config.RedisConfig) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return &Redis{}, err
	}

	return &Redis{
		Client: client,
	}, nil
}

func (r *Redis) Close() error {
	return r.Client.Close()
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}
