package main

import (
	"fmt"
	"github.com/abolfazlalz/goasali/kit/cache"
	"time"
)

func main() {
	fmt.Println("Test redis cache connection")
	cache := cache.Redis[string]{}
	cache.Set("hello", "world", time.Second*10)
	fmt.Println(cache.Get("hello"))
}
