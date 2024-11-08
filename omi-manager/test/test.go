package main

import (
	"github.com/go-redis/redis/v8"
	omimanager "github.com/stormi-li/omi/omi-manager"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	c:=omimanager.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})
	c.Start("118.25.196.166:8888")
}
