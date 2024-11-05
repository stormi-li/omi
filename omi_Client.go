package omi

import (
	"context"

	"github.com/go-redis/redis/v8"
	omipc "github.com/stormi-li/omi/omi-ipc"
)

type Client struct {
	redisClient *redis.Client
	omipcClient *omipc.Client
	namespace   string
	serverType  ServerType
}

func (c *Client) NewRegister(serverName string, address string) *Register {
	return &Register{
		redisClient:      c.redisClient,
		omipcClient:      c.omipcClient,
		namespace:        c.namespace,
		serverName:       serverName,
		ctx:              context.Background(),
		address:          address,
		redisChannelName: serverName + NamespaceSeparator + address,
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

func (c *Client) NewConsumer(channel string, address string) *Consumer {
	if c.serverType != MQ {
		panic("server type must be mq")
	}
	return &Consumer{
		omiClient: c,
		channel:   channel,
		address:   address,
		Register:  c.NewRegister(channel, address),
	}
}

func (c *Client) NewProducer(channel string) *Producer {
	if c.serverType != MQ {
		panic("server type must be mq")
	}
	producer := Producer{
		omiClient: c,
		channel:   channel,
	}
	go producer.listen()
	return &producer
}
