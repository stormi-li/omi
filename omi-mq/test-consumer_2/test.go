package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	omique "github.com/stormi-li/omi/omi-mq"
)

func main() {
	consumer()
}

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func consumer() {
	c := omique.NewClient(&redis.Options{Addr: redisAddr, Password: password})
	consumer := c.NewConsumer("consumer_1", 1)
	consumer.AddHandler(func(message []byte) {
		fmt.Println(string(message))
	})
	consumer.Listen("118.25.196.166:5556")
}
