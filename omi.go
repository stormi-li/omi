package omi

import (
	"github.com/go-redis/redis/v8"
	omiclient "github.com/stormi-li/omi/omi-client"
)

func NewServerClient(redisClient *redis.Client, namespace string) *omiclient.Client {
	return omiclient.NewClient(redisClient, namespace, server, const_serverPrefix)
}

func NewMQClient(redisClient *redis.Client, namespace string) *omiclient.MQClient {
	return omiclient.NewMQClient(redisClient, namespace, mq, const_mqPrefix)
}

func NewConfigClient(redisClient *redis.Client, namespace string) *omiclient.Client {
	return omiclient.NewClient(redisClient, namespace, config, const_configPrefix)
}
