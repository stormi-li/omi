package omi

import (
	"github.com/go-redis/redis/v8"
	omipc "github.com/stormi-li/omi/om-ipc"
)

type Client struct {
	redisClient *redis.Client
	omipcClient *omipc.Client
	namespace   string
	serverType  ServerType
}

func NewClient(redisClient *redis.Client, namespace string, serverType ServerType) *Client {
	prefix := ""
	if serverType == Server {
		prefix = const_serverPrefix
	}
	if serverType == MQ {
		prefix = const_mqPrefix
	}
	if serverType == Config {
		prefix = const_configPrefix
	}
	return &Client{
		omipcClient: omipc.NewClient(redisClient, namespace),
		redisClient: redisClient,
		namespace:   namespace + const_separator + prefix,
		serverType:  serverType,
	}
}

func (c *Client) GetOmipc() *omipc.Client {
	return c.omipcClient
}

func (c *Client) NewRegister(serverName string, address string) *Register {
	return newRegister(c.redisClient, c.omipcClient, c.namespace, serverName, address)
}

func (c *Client) NewSearcher() *Searcher {
	return newSearcher(c.redisClient, c.omipcClient, c.namespace)
}

func (c *Client) NewConsumer(channel string, address string) *Consumer {
	if c.serverType != MQ {
		panic("server type must be mq")
	}
	return newConsumer(c, channel, address)
}

func (c *Client) NewProducer(channel string) *Producer {
	if c.serverType != MQ {
		panic("server type must be mq")
	}
	return newProducer(c, channel)
}

func NewManager(redisClient *redis.Client, namespace string) *Manager {
	return newManager(redisClient, namespace)
}
