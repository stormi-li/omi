package omiweb

import (
	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi"
)

func NewClient(redisClient *redis.Client, namespace string, serverName, address string) *Client {
	omiClient := omi.NewClient(redisClient, namespace, omi.Server)
	return &Client{
		router:      newRouter(omiClient.NewSearcher()),
		redisClient: redisClient,
		omiClient:   omiClient,
		serverName:  serverName,
		address:     address,
		namespace:   namespace,
	}
}
