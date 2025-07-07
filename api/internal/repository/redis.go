package repository

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
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

	if cfg.MasterName != "" && len(cfg.Addrs) > 0 && cfg.Addrs[0] != "" {
		client, err = newFailoverClient(cfg)
	}

	if err != nil {
		log.Println("Failed to connect to Redis:", err.Error())
		return nil, err
	}

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	log.Println("Redis connection OK")

	return &Redis{Client: client}, nil
}

func newClient(cfg config.RedisConfig) (*redis.Client, error) {
	log.Println("Creating Redis client")
	options := &redis.Options{
		Addr: cfg.Addr,
	}

	return redis.NewClient(options), nil
}

func newFailoverClient(cfg config.RedisConfig) (*redis.Client, error) {
	log.Println("Creating Redis failover client")
	options := &redis.FailoverOptions{
		MasterName:       cfg.MasterName,
		Username:         cfg.Username,
		Password:         cfg.Password,
		SentinelUsername: cfg.FailoverUsername,
		SentinelPassword: cfg.FailoverPassword,
		SentinelAddrs:    cfg.Addrs,
		DB:               0,
	}

	if cfg.TLSEnabled {
		log.Println("Using TLS to connect to Redis")

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

func (r *Redis) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *Redis) Del(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (c *Redis) Incr(ctx context.Context, key string, expiration time.Duration) error {
	err := c.Client.Incr(ctx, key).Err()
	if err != nil {
		return err
	}

	if expiration > 0 {
		return c.Client.Expire(ctx, key, expiration).Err()
	}

	return nil
}
