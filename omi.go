package omi

import (
	"github.com/go-redis/redis/v8"
	omipc "github.com/stormi-li/omi/omi-ipc"
)

func NewClient(redisClient *redis.Client, namespace string, serverType ServerType) *Client {
	prefix := ""
	if serverType == Server {
		prefix = const_serverPrefix
	}
	if serverType == MQ {
		prefix = const_mqPrefix
	}
	if serverType == Config {
		prefix = const_configPrefix
	}
	return &Client{
		omipcClient: omipc.NewClient(redisClient, namespace),
		redisClient: redisClient,
		namespace:   namespace + const_separator + prefix,
		serverType:  serverType,
	}
}
