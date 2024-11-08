package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	omiC := omi.NewConfigClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})
	web := omi.NewWebClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})
	web.NewRegister("mysql", "118.25.196.166:3306").Start(3, map[string]string{})
	omiC.NewRegister("mysql", "118.25.196.166:3306").Start(3, map[string]string{})
	omiC.NewRegister("redis", "118.25.196.166:6379").Start(3, map[string]string{})
	omiC.NewRegister("redis", "118.25.196.166:6378").Start(4, map[string]string{})
	select {}
}
