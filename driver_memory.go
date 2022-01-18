package cache

import (
	"sync"
	"time"
)

type Item struct {
	Exp int64
	Val interface{}
}

func (item Item) Expired() bool {
	return item.Exp > 1 && item.Exp < time.Now().Unix()
}

type MemoryCache struct {
	lock      sync.RWMutex
	caches    map[string]*Item
	CacheSize int
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		caches: make(map[string]*Item),
	}
}

func (c *MemoryCache) Has(key string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.get(key) != nil
}

func (c *MemoryCache) Get(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.get(key)
}

func (c *MemoryCache) get(key string) interface{} {
	if item, ok := c.caches[key]; ok {
		if item.Expired() {
			_ = c.del(key)
			return nil
		}

		return item.Val
	}

	return nil
}

func (c *MemoryCache) Set(key string, val interface{}, ttl time.Duration) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.set(key, val, ttl)
}

func (c *MemoryCache) set(key string, val interface{}, ttl time.Duration) (err error) {
	item := &Item{Val: val}
	if ttl > 0 {
		item.Exp = time.Now().Unix() + int64(ttl/time.Second)
	}

	c.caches[key] = item
	return
}

func (c *MemoryCache) Del(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.del(key)
}

func (c *MemoryCache) del(key string) error {
	if _, ok := c.caches[key]; ok {
		delete(c.caches, key)
	}

	return nil
}

func (c *MemoryCache) SetMulti(values map[string]interface{}, ttl time.Duration) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	for key, val := range values {
		if err = c.set(key, val, ttl); err != nil {
			return
		}
	}
	return
}

func (c *MemoryCache) DelMulti(keys []string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, key := range keys {
		_ = c.del(key)
	}

	return nil
}

func (c *MemoryCache) Count() int {
	return len(c.caches)
}
