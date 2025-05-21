package lru

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"time"

	golru "github.com/hashicorp/golang-lru/v2"
)

const Limit = 2000

var ErrNotFound = errors.New("not found")

type Client struct {
	cache *golru.Cache[string, *Item]
}

func New() (*Client, error) {
	c, err := golru.New[string, *Item](Limit)
	if err != nil {
		return nil, err
	}
	return &Client{cache: c}, nil
}

func (c *Client) Get(ctx context.Context, key string, out any) error {
	if v, ok := c.cache.Get(key); ok && v.TTL.After(time.Now()) {
		return v.Scan(out)
	}

	return ErrNotFound
}

func (c *Client) Set(ctx context.Context, key string, in any, ttl time.Duration) error {
	data, err := newItem(in, ttl)
	if err != nil {
		return err
	}
	c.cache.Add(key, data)
	return nil
}

func (c *Client) DelAll(ctx context.Context) error {
	c.cache.Purge()
	return nil
}

func (c *Client) Del(ctx context.Context, key string) error {
	c.cache.Remove(key)
	return nil
}

type Item struct {
	Data []byte
	TTL  time.Time
}

func (c *Item) Scan(out any) error {
	buf := bytes.NewReader(c.Data)
	decoder := gob.NewDecoder(buf)
	return decoder.Decode(out)
}

func newItem(data any, ttl time.Duration) (*Item, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	return &Item{
		Data: buf.Bytes(),
		TTL:  time.Now().Add(ttl),
	}, nil
}
