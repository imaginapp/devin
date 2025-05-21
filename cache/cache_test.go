package cache

import (
	"context"
	"testing"
	"time"

	"github.com/imaginapp/devin/lru"
	"github.com/imaginapp/proto/go/gen/imagin/external/message/v1"
	"github.com/stretchr/testify/assert"
)

func TestCacheItem(t *testing.T) {
	assert := assert.New(t)

	lru, err := lru.New()
	assert.NoError(err)
	c := New(lru)

	key := "accountData"
	data := &message.Account{
		Id: "testabd123",
	}
	err = c.Set(context.Background(), key, data, time.Second*10)
	assert.NoError(err)
	outData := &message.Account{}
	err = c.Get(context.Background(), key, outData)
	assert.NoError(err)
}

func TestCacheItemNotFound(t *testing.T) {
	assert := assert.New(t)

	lru, err := lru.New()
	assert.NoError(err)
	c := New(lru)

	key := "accountData"
	var outData any
	err = c.Get(context.Background(), key, &outData)
	assert.ErrorIs(err, ErrCacheNotFound)
}
