package omique

import (
	"github.com/go-redis/redis/v8"
	omiclient "github.com/stormi-li/omi/omi-client"
)

type Client struct {
	omiClient *omiclient.Client
}

func newClient(redisClient *redis.Client, namespace string, serverType string, prefix string) *Client {
	return &Client{omiClient: omiclient.NewClient(redisClient, namespace, serverType, prefix)}
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
		omiClient: c.omiClient,
		channel:   channel,
	}
	producer.listen()
	return &producer
}
