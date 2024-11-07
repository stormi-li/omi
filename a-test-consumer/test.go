package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	omique "github.com/stormi-li/omi/omi-mq"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})
	client := omique.NewClient(redisClient, "omi-chat")
	consumer := client.NewConsumer("consumer_test")
	consumer.ListenOnMain("118.25.196.166:4444", func(message []byte) {
		fmt.Println(string(message))
	})
}
