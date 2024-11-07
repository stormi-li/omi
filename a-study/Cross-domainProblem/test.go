package main

import (
	"github.com/go-redis/redis/v8"
	omiweb "github.com/stormi-li/omi/omi-web"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})
	omiwebC := omiweb.NewClient(redisClient, "omi-chat")
	omiwebC.GenerateTemplate()
	omiwebC.Listen("118.25.196.166:9911")

}
