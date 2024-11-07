package omiweb

import (
	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi"
)

func NewClient(redisClient *redis.Client, namespace string) *Client {
	omiClient := omi.NewServerClient(redisClient, namespace)
	return &Client{
		redisClient: redisClient,
		omiClient:   omiClient,
		namespace:   namespace,
	}
}
