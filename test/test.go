package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi"
)

func main() {
	omi.NewWebClient(&redis.Options{})
}
