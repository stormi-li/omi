package omiweb

import (
	"github.com/go-redis/redis/v8"
	omiclient "github.com/stormi-li/omi/omi-client"
)

func NewClient(redisClient *redis.Client, omiClient *omiclient.Client) *Client {
	return &Client{
		redisClient: redisClient,
		omiClient:   omiClient,
	}
}
