package omiclient

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	redisClient *redis.Client
	namespace   string
	serverType  string
}

func NewClient(redisClient *redis.Client, serverType string, prefix string) *Client {
	return &Client{
		redisClient: redisClient,
		namespace:   prefix,
		serverType:  serverType,
	}
}

func (c *Client) NewRegister(serverName string, weight int) *Register {
	return &Register{
		redisClient: c.redisClient,
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
		namespace:   c.namespace,
		ctx:         context.Background(),
	}
}
