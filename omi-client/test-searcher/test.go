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
	address, data := searcher.SearchByLoadBalancing("mysql")
	fmt.Println(address, data)
	fmt.Println(searcher.IsAlive("mysql", address))
}
