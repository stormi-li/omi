package omi

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
	omipc "github.com/stormi-li/omi/om-ipc"
)

type Client struct {
	redisClient *redis.Client
	omipcClient *omipc.Client
	namespace   string
	serverType  ServerType
}

func (c *Client) GetOmipc() *omipc.Client {
	return c.omipcClient
}

func (c *Client) NewRegister(serverName string, address string) *Register {
	return &Register{
		redisClient: c.redisClient,
		omipcClient: c.omipcClient,
		namespace:   c.namespace,
		serverName:  serverName,
		ctx:         context.Background(),
		address:     address,
		channel:     serverName + const_separator + address,
		CloseSignal: make(chan struct{}, 1),
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
		omiClient:   c,
		channel:     channel,
		address:     address,
		messageChan: make(chan []byte, 1000),
		buffer:      [][]byte{},
		bufferLock:  sync.Mutex{},
		Register:    c.NewRegister(channel, address),
	}
}

func (c *Client) NewProducer(channel string) *Producer {
	if c.serverType != MQ {
		panic("server type must be mq")
	}
	producer := Producer{
		omiClient:  c,
		maxRetries: 10,
		channel:    channel,
	}
	go producer.omiClient.NewSearcher().Listen(producer.channel, func(addr string, data map[string]string) {
		producer.address = addr
		producer.connect()
	})
	return &producer
}
