package cache

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"time"
)

const (
	NoExpiration      = -1
	DefaultExpiration = 0
)

var (
	ErrorKeyNotFound = errors.New("key not found")
)

type Item struct {
	Value      interface{}
	Expiration int64
}

type Cache struct {
	defaultExpiration time.Duration
	ctx               context.Context
	inMem             map[string]Item
	mut               sync.RWMutex
	janitor           *janitor
}

func NewCache(ctx context.Context, defaultExpiration time.Duration, janitorInterval time.Duration) *Cache {
	if janitorInterval == 0 {
		janitorInterval = 2 * time.Minute
	}

	c := &Cache{
		ctx:   ctx,
		inMem: make(map[string]Item),
	}

	j := newJanitor(janitorInterval)
	j.Run(c)
	runtime.SetFinalizer(c, stopJanitor)

	return c
}

func (c *Cache) Get(key string) (interface{}, error) {
	c.mut.RLock()
	defer c.mut.RUnlock()
	if val, ok := c.inMem[key]; ok {
		return val.Value, nil
	}
	return nil, ErrorKeyNotFound
}

func (c *Cache) SetDefault(key string, val interface{}) error {
	return c.Set(key, val, c.defaultExpiration)
}

func (c *Cache) Set(key string, val interface{}, d time.Duration) error {
	var exp int64

	if d == DefaultExpiration {
		d = c.defaultExpiration
	}

	if d > 0 {
		exp = time.Now().Add(d).UnixNano()
	}

	c.mut.Lock()
	c.inMem[key] = Item{
		Value:      val,
		Expiration: exp,
	}
	c.mut.Unlock()

	return nil
}

func (c *Cache) Delete(key string) error {
	c.mut.Lock()
	defer c.mut.Unlock()
	delete(c.inMem, key)
	return nil
}

func (c *Cache) DeleteExpired() {
	now := time.Now().UnixNano()
	c.mut.Lock()
	for k, v := range c.inMem {
		if v.Expiration > 0 && now > v.Expiration {
			delete(c.inMem, k)
		}
	}
	c.mut.Unlock()
}
