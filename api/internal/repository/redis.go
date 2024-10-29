package repository

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"ivpn.net/email/api/config"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis(cfg config.RedisConfig) (*Redis, error) {
	var client *redis.Client
	client, err := newClient(cfg)

	if cfg.MasterName != "" && len(cfg.Addrs) > 0 {
		client, err = newFailoverClient(cfg)
	}

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &Redis{Client: client}, nil
}

func newClient(cfg config.RedisConfig) (*redis.Client, error) {
	options := &redis.Options{
		Addr: cfg.Addr,
	}

	return redis.NewClient(options), nil
}

func newFailoverClient(cfg config.RedisConfig) (*redis.Client, error) {
	options := &redis.FailoverOptions{
		MasterName:    cfg.MasterName,
		SentinelAddrs: cfg.Addrs,
		Password:      "",
		DB:            0,
	}

	if cfg.TLSEnabled {
		cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %v", err)
		}

		caCert, err := os.ReadFile(cfg.CACertFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load CA certificate: %v", err)
		}

		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			return nil, fmt.Errorf("failed to append CA certificate")
		}

		options.TLSConfig = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: cfg.TLSInsecureSkipVerify, // Only for testing, use false in production
		}
	}

	return redis.NewFailoverClient(options), nil
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

func (r *Redis) Del(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
