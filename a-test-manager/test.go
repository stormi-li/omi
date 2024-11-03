package main

import (
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
	manager := omi.NewManager(redisClient, "omi-namespace")
	manager.Start("omi-manager", "118.25.196.166:7766")
}
