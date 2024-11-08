package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})
	c := omi.NewServerClient(redisClient, "omi-namespace")
	searcher := c.NewSearcher()
	for _, val := range searcher.SearchRunningServers("user_server") {
		fmt.Println(val)
	}
}
