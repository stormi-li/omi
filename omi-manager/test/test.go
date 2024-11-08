package main

import (
	"github.com/go-redis/redis/v8"
	omimanager "github.com/stormi-li/omi/omi-manager"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	omimanager.Start(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	}, "118.25.196.166:8080")
}
