package storage_redis

import (
	"context"
	"mess/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type CashInterface interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, time time.Duration) error
	Del(ctx context.Context, key string) error
}

type Client struct {
	client *redis.Client
}

func NewClient(cfg *config.Config) *Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis,
		Password: "",
		DB:       0,
	})
	return &Client{client: client}
}

// Get result from Redis
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Set result to Redis
func (c *Client) Set(ctx context.Context, key string, value string, time time.Duration) error {
	return c.client.Set(ctx, key, value, time).Err()
}

// Delete result from Redis
func (c *Client) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
