package omique

import (
	"github.com/go-redis/redis/v8"
	omiclient "github.com/stormi-li/omi/omi-client"
)

func NewClient(redisClient *redis.Client, namespace string) *Client {
	return newClient(redisClient, namespace, omiclient.Config, omiclient.Prefix_Config)
}
