package main

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	omiClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})

	go func() {
		client := omi.NewServerClient(omiClient, "omi-namespace")
		register := client.NewRegister("user_server", "1223213:2222")
		register.StartOnBackup(map[string]string{"message": "user_server"})
	}()
	go func() {
		client := omi.NewServerClient(omiClient, "omi-namespace")
		register := client.NewRegister("user_server", "1223213:2221")
		register.StartOnMain(map[string]string{"message": "user_server"})
	}()
	go func() {
		client := omi.NewServerClient(omiClient, "omi-namespace")
		register := client.NewRegister("user_server", "1223213:2220")
		go func() {
			time.Sleep(1 * time.Second)
			register.ToStop()
		}()
		register.StartOnMain(map[string]string{"message": "user_server"})
	}()
	go func() {
		client := omi.NewServerClient(omiClient, "omi-namespace")
		register := client.NewRegister("order_server", "1223213:1111")
		register.StartOnBackup(map[string]string{"message": "order_server"})
	}()
	go func() {
		client := omi.NewServerClient(omiClient, "omi-namespace")
		register := client.NewRegister("order_server", "1223213:1112")
		register.StartOnMain(map[string]string{"message": "order_server"})
	}()
	go func() {
		client := omi.NewServerClient(omiClient, "omi-namespace")
		register := client.NewRegister("order_server", "1223213:1113")
		go func() {
			time.Sleep(1 * time.Second)
			register.ToStop()
		}()
		register.StartOnMain(map[string]string{"message": "order_server"})
	}()
	go func() {
		client := omi.NewConfigClient(omiClient, "omi-namespace")
		register := client.NewRegister("mysql", "1223213:3306")
		register.StartOnBackup(map[string]string{"message": "mysql"})
	}()
	go func() {
		client := omi.NewConfigClient(omiClient, "omi-namespace")
		register := client.NewRegister("mysql", "1223213:3307")
		register.StartOnMain(map[string]string{"message": "mysql"})
	}()
	go func() {
		client := omi.NewConfigClient(omiClient, "omi-namespace")
		register := client.NewRegister("mysql", "1223213:3308")
		go func() {
			time.Sleep(1 * time.Second)
			register.ToStop()
		}()
		register.StartOnMain(map[string]string{"message": "mysql"})
	}()
	go func() {
		client := omi.NewConfigClient(omiClient, "omi-namespace")
		register := client.NewRegister("redis-config-非常重要------------------------------------------------------", "1223213:6379")
		register.StartOnBackup(map[string]string{"message": "redis"})
	}()
	go func() {
		client := omi.NewConfigClient(omiClient, "omi-namespace")
		register := client.NewRegister("redis-config-非常重要------------------------------------------------------", "1223213:6378")
		register.StartOnMain(map[string]string{"message": "redis"})
	}()
	go func() {
		client := omi.NewConfigClient(omiClient, "omi-namespace")
		register := client.NewRegister("redis-config-非常重要------------------------------------------------------", "1223213:6377")
		go func() {
			time.Sleep(1 * time.Second)
			register.ToStop()
		}()
		register.StartOnMain(map[string]string{"message": "redis"})
	}()

	select {}
}
