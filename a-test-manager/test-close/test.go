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
	client := omi.NewClient(redisClient, "omi-namespace", omi.Server)
	rg := client.NewRegister("omi-web-manager", "118.25.196.166:7766")
	rg.Close()
}
