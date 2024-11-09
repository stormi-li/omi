package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	searcher := omi.NewConfigClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	}).NewSearcher()
	name, data := searcher.SearchOneByWeight("redis")
	fmt.Println(name, data)
}
