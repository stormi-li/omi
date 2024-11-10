package omiweb

import (
	"github.com/go-redis/redis/v8"
	omiclient "github.com/stormi-li/omi/omi-client"
)

func NewClient(redisClient *redis.Client, omiWebClient, omiServerClient *omiclient.Client) *Client {
	return &Client{
		redisClient:     redisClient,
		omiWebClient:    omiWebClient,
		omiServerClient: omiServerClient,
	}
}

func DisableLog() {
	log_cache = false
}
