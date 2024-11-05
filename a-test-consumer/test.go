package main

import (
	"fmt"

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
	client := omi.NewClient(redisClient, "omi-namespace", omi.MQ)
	consumer := client.NewConsumer("consumer_test", "118.25.196.166:4443")
	consumer.StartOnMain(1000000, func(message []byte) {
		fmt.Println(string(message))
	})

}
