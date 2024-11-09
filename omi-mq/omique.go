package omique

import (
	"github.com/go-redis/redis/v8"
	omiclient "github.com/stormi-li/omi/omi-client"
)

func NewClient(opts *redis.Options) *Client {
	return newClient(redis.NewClient(opts), omiclient.Config, omiclient.Prefix_Config)
}
