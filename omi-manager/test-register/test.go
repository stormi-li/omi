package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	omiC := omi.NewConfigManager(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})
	r := omiC.NewRegister("mysql", 1)
	r.Data["username"] = "root"
	r.Data["database"] = "USER"
	r.Data["password"] = "12982397StrongPassw0rd"
	r.Register("118.25.196.166:3933")
	r = omiC.NewRegister("redis", 1)
	r.Register("118.25.196.166:6379")
	r = omiC.NewRegister("redis1", 1)
	r.Register("118.25.196.166:6379")
	select {}
}
