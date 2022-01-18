package cache_test

import (
	"fmt"
	"github.com/liyaojian/cache"
	"time"
)

func ExampleFileCache() {
	c := cache.NewFileCache("./testdata")
	key := "name"

	// set
	c.Set(key, "cache value", time.Minute)
	fmt.Println(c.Has(key), c.Count())

	// get
	val := c.Get(key)
	fmt.Println(val)

	// del
	c.Del(key)
	fmt.Println(c.Has(key), c.Count())

	// Output:
	// true 1
	// cache value
	// false 0
}
