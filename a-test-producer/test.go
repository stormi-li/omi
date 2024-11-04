package main

import (
	"strconv"
	"time"

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
	producer := client.NewProducer("consumer_test")
	for i := 0; i < 10000; i++ {
		producer.Publish([]byte("omi" + strconv.Itoa(i)))
		time.Sleep(100 * time.Millisecond)
	}
}
