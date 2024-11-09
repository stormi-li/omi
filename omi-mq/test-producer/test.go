package main

import (
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	omique "github.com/stormi-li/omi/omi-mq"
)

func main() {
	producer()
}

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func producer() {
	c := omique.NewClient(&redis.Options{Addr: redisAddr, Password: password})
	producer := c.NewProducer("consumer_1")
	for i := 0; i < 10000; i++ {
		producer.Publish([]byte("message" + strconv.Itoa(i)))
		time.Sleep(50 * time.Millisecond)
	}
}
