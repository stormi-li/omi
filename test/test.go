package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	web := omi.NewWebClient(&redis.Options{Addr: redisAddr, Password: password})
	web.GenerateTemplate()
	webServer := web.NewWebServer("web demo", 1)
	webServer.Listen("118.25.196.166:8848")
}
