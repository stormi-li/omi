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

	client := omi.NewConfigClient(redisClient, "omi-chat")
	register := client.NewRegister("mysql", "118.25.196.166:3933")
	register.StartOnBackup(map[string]string{"username": "root", "password": "12982397StrongPassw0rd", "database": "USER"})

	select {}
}
