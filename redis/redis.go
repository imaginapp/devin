package redis

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

var ErrNotFound = errors.New("not found")

type Config goredis.Options

type Client struct {
	store *goredis.Client
}

func New(config *Config) (*Client, error) {
	if config == nil {
		return nil, errors.New("redis config required")
	}
	redisOptions := goredis.Options(*config)

	client := goredis.NewClient(&redisOptions)

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("could not create cache: %w", err)
	}

	return &Client{store: client}, nil
}

func NewFromEnv(DB int) (*Client, error) {
	redisAddress := os.Getenv("REDIS_ADDRESS")
	if redisAddress == "" {
		err := errors.New("REDIS_ADDRESS not found in ENV")
		return nil, err
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		err := errors.New("REDIS_PORT not found in ENV")
		return nil, err
	}

	return New(&Config{
		Addr: fmt.Sprintf("%s:%s", redisAddress, redisPort),
		DB:   1,
	})
}

func (c *Client) Set(ctx context.Context, key string, v any, ttl time.Duration) error {
	return c.store.Set(ctx, key, v, ttl).Err()
}

func (c *Client) Get(ctx context.Context, key string, out any) error {
	err := c.store.Get(ctx, key).Scan(out)
	if err != nil {
		if err == goredis.Nil {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (c *Client) DelAll(ctx context.Context) error {
	return c.store.FlushAll(ctx).Err()
}

func (c *Client) Del(ctx context.Context, key string) error {
	return c.store.Del(ctx, key).Err()
}
