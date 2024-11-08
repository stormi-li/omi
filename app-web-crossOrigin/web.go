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
	c := omiweb.NewClient(redisClient, "omi-chat")
	c.GenerateTemplate()
	ws := c.NewWebServer("跨域界面")
	ws.Listen("118.25.196.166:8083")
}
