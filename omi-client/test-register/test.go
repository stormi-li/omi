package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	omiC := omi.NewConfigClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})
	omiC.NewRegister("redis", 2).Register("118.25.196.166:6378")
	select {}
}
