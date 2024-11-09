package omiclient

import (
	"context"

	"github.com/go-redis/redis/v8"
	omipc "github.com/stormi-li/omi/omi-ipc"
)

type Client struct {
	redisClient *redis.Client
	omipcClient *omipc.Client
	namespace   string
	serverType  string
}

func NewClient(redisClient *redis.Client, serverType string, prefix string) *Client {
	return &Client{
		omipcClient: omipc.NewClient(redisClient),
		redisClient: redisClient,
		namespace:   prefix,
		serverType:  serverType,
	}
}

func (c *Client) NewRegister(serverName string, weight int) *Register {
	return &Register{
		redisClient: c.redisClient,
		omipcClient: c.omipcClient,
		serverName:  serverName,
		weight:      weight,
		Data:        map[string]string{},
		namespace:   c.namespace,
		ctx:         context.Background(),
	}
}

func (c *Client) NewSearcher() *Searcher {
	return &Searcher{
		redisClient: c.redisClient,
		omipcClient: c.omipcClient,
		namespace:   c.namespace,
		ctx:         context.Background(),
	}
}
