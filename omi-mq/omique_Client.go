package omique

import (
	"github.com/go-redis/redis/v8"
	omiclient "github.com/stormi-li/omi/omi-manager"
)

type Client struct {
	omiClient *omiclient.Client
}

func newClient(redisClient *redis.Client, serverType string, prefix string) *Client {
	return &Client{omiClient: omiclient.NewClient(redisClient, serverType, prefix)}
}

func (c *Client) NewConsumer(channel string, weight int) *Consumer {
	return &Consumer{
		omiClient:   c.omiClient,
		channel:     channel,
		weight:      weight,
		messageChan: make(chan []byte, 1000000),
	}
}

func (c *Client) NewProducer(channel string) *Producer {
	producer := Producer{
		searcher: c.omiClient.NewSearcher(),
		channel:  channel,
	}
	return &producer
}
