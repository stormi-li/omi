package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	omi.NewConfigClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	}).NewRegister("mysql", "118.25.196.166:3933").Start(1, map[string]string{"username": "root", "database": "USER", "password": "12982397StrongPassw0rd"})
	select{}
}
