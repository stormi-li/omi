package omipc

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func NewClient(redisClient *redis.Client, namespace string) *Client {
	return &Client{
		redisClient: redisClient,
		namespace:   namespace + NamespaceSeparator,
		ctx:         context.Background(),
	}
}
