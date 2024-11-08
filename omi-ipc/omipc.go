package omipc

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func NewClient(redisClient *redis.Client) *Client {
	return &Client{
		redisClient: redisClient,
		ctx:         context.Background(),
	}
}
