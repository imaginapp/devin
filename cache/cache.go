package cache

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
)

type Store interface {
	Get(ctx context.Context, key string, value any) error
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
}

type Marshaler interface {
	Marshal(v any) ([]byte, error)
}

type Client struct {
	stores []Store
}

const DefaultCacheSec = 360
const MaxCacheSec = 3600

var ErrNoStores = errors.New("no cache stores configured")
var ErrCacheNotFound = errors.New("cache not found")

func New(stores ...Store) *Client {
	return &Client{
		stores: stores,
	}
}

func (c *Client) Set(ctx context.Context, key string, data any, ttl time.Duration) error {
	if len(c.stores) == 0 {
		return nil
	}
	for _, store := range c.stores {
		err := store.Set(ctx, key, data, ttl)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) Get(ctx context.Context, key string, out any) error {
	if len(c.stores) == 0 {
		return ErrNoStores
	}

	for idx, store := range c.stores {
		err := store.Get(ctx, key, out)
		if err == nil {
			if idx > 0 {
				// promote to cache 0
				if pErr := c.stores[0].Set(ctx, key, out, time.Second*DefaultCacheSec); pErr != nil {
					log.Debug().Err(pErr).Msg("failed to promote cached item")
				}
			}
			// if no error then its probably found the item
			return nil
		}
		log.Debug().Err(err).Msg("failed to get cached item")
	}

	return ErrCacheNotFound
}
