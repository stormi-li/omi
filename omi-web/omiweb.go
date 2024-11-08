package omiweb

import (
	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi"
)

func NewClient(opts *redis.Options) *Client {
	omiClient := omi.NewServerClient(opts)
	return &Client{
		redisClient: redis.NewClient(opts),
		omiClient:   omiClient,
	}
}
