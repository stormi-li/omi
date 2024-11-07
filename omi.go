package omi

import (
	"github.com/go-redis/redis/v8"
	omiclient "github.com/stormi-li/omi/omi-client"
)

func NewServerClient(redisClient *redis.Client, namespace string) *omiclient.Client {
	return omiclient.NewClient(redisClient, namespace, omiclient.Server, omiclient.Prefix_Server)
}

func NewWebClient(redisClient *redis.Client, namespace string) *omiclient.Client {
	return omiclient.NewClient(redisClient, namespace, omiclient.Web, omiclient.Prefix_Web)
}

func NewConfigClient(redisClient *redis.Client, namespace string) *omiclient.Client {
	return omiclient.NewClient(redisClient, namespace, omiclient.Config, omiclient.Prefix_Config)
}
