package main

import (
	"embed"

	"github.com/go-redis/redis/v8"
	omiweb "github.com/stormi-li/omi/omi-web"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

//go:embed src/*
var embedSource embed.FS

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})
	omiweb := omiweb.NewClient(redisClient, "omi-namespace", "web-server", "118.25.196.166:7788")
	omiweb.Start(embedSource)
}
