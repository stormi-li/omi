package omi

import (
	"github.com/go-redis/redis/v8"
	omiclient "github.com/stormi-li/omi/omi-client"
)

func NewServerClient(opts *redis.Options) *omiclient.Client {
	return omiclient.NewClient(redis.NewClient(opts), omiclient.Server, omiclient.Prefix_Server)
}

func NewWebClient(opts *redis.Options) *omiclient.Client {
	return omiclient.NewClient(redis.NewClient(opts), omiclient.Web, omiclient.Prefix_Web)
}

func NewConfigClient(opts *redis.Options) *omiclient.Client {
	return omiclient.NewClient(redis.NewClient(opts), omiclient.Config, omiclient.Prefix_Config)
}
