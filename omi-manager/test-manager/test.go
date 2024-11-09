package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	c := omi.NewManager(&redis.Options{Addr: redisAddr, Password: password})
	c.Listen("118.25.196.166:9999")
}
