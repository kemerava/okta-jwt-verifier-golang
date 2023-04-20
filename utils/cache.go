package utils

import (
	"net/http"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

// Cacher is a read-only cache interface.
//
// Get returns the value associated with the given key.
type Cacher interface {
	Get(string, *http.Client) (interface{}, error)
}

type defaultCache struct {
	cache  *cache.Cache
	lookup func(string, *http.Client) (interface{}, error)
	mutex  *sync.Mutex
}

func (c *defaultCache) Get(key string, client *http.Client) (interface{}, error) {
	if value, found := c.cache.Get(key); found {
		return value, nil
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// once lock, check the cache again because there could be
	// another thread that has update the keys during the last check
	if value, found := c.cache.Get(key); found {
		return value, nil
	}

	value, err := c.lookup(key, client)
	if err != nil {
		return nil, err
	}

	c.cache.SetDefault(key, value)
	return value, nil
}

// defaultCache implements the Cacher interface
var _ Cacher = (*defaultCache)(nil)

func NewDefaultCache(lookup func(string, *http.Client) (interface{}, error), timeout, cleanup time.Duration) (Cacher, error) {
	return &defaultCache{
		cache:  cache.New(timeout, cleanup),
		lookup: lookup,
		mutex:  &sync.Mutex{},
	}, nil
}
