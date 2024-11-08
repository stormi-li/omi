package omique

import (
	"github.com/go-redis/redis/v8"
	omiclient "github.com/stormi-li/omi/omi-client"
)

type Client struct {
	omiClient *omiclient.Client
}

func newClient(redisClient *redis.Client, serverType string, prefix string) *Client {
	return &Client{omiClient: omiclient.NewClient(redisClient, serverType, prefix)}
}

func (c *Client) NewConsumer(channel string) *Consumer {
	return &Consumer{
		omiClient:   c.omiClient,
		channel:     channel,
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