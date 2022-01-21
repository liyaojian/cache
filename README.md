File Cache

Usage:
```go
package main

import (
	"fmt"
	"github.com/liyaojian/cache"
	"time"
)

func main() {
	c := cache.NewFileCache("./cache")
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
}
```