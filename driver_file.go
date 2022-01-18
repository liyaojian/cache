package cache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/liyaojian/cache/utils"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type FileCache struct {
	MemoryCache
	Dir    string
	Prefix string
}

func NewFileCache(dir string, prefix ...string) *FileCache {
	utils.MustExist(dir)
	c := &FileCache{
		MemoryCache: MemoryCache{caches: make(map[string]*Item)},
		Dir:         dir,
	}
	if len(prefix) > 0 {
		c.Prefix = prefix[0]
	}
	return c
}

func (c *FileCache) path(key string) string {
	h := md5.New()
	h.Write([]byte(key))
	str := hex.EncodeToString(h.Sum(nil))

	return strings.Join([]string{
		c.Dir,
		c.Prefix + str + ".data",
	}, "/")
}

func (c *FileCache) Has(key string) bool {
	return c.get(key) != nil
}

func (c *FileCache) Get(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.get(key)
}

func (c *FileCache) get(key string) interface{} {
	// read cache from memory
	if val := c.MemoryCache.get(key); val != nil {
		return val
	}

	// read cache from file
	bs, err := ioutil.ReadFile(c.path(key))
	if err != nil {
		return nil
	}

	item := &Item{}
	if err = json.Unmarshal(bs, item); err != nil {
		return nil
	}

	// check expired
	if item.Expired() {
		_ = c.del(key)
		return nil
	}

	c.caches[key] = item // save to memory.
	return item.Val
}

func (c *FileCache) set(key string, val interface{}, ttl time.Duration) (err error) {
	err = c.MemoryCache.set(key, val, ttl)
	if err != nil {
		return
	}

	// cache item data to file
	bs, err := json.Marshal(c.caches[key])
	if err != nil {
		return
	}

	file := c.path(key)
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(bs)
	return
}

func (c *FileCache) Set(key string, val interface{}, ttl time.Duration) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.set(key, val, ttl)
}

func (c *FileCache) Del(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.del(key)
}

func (c *FileCache) del(key string) error {
	if err := c.MemoryCache.del(key); err != nil {
		return err
	}

	file := c.path(key)
	if utils.PathExists(file) {
		return os.Remove(file)
	}

	return nil
}

func (c *FileCache) GetMulti(keys []string) map[string]interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	data := make(map[string]interface{}, len(keys))
	for _, key := range keys {
		data[key] = c.get(key)
	}

	return data
}

func (c *FileCache) SetMulti(values map[string]interface{}, ttl time.Duration) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for key, val := range values {
		if err = c.set(key, val, ttl); err != nil {
			return
		}
	}
	return
}

func (c *FileCache) DelMulti(keys []string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, key := range keys {
		_ = c.del(key)
	}
	return nil
}
