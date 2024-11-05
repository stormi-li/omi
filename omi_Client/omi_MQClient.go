package omiclient

import "github.com/go-redis/redis/v8"

type MQClient struct {
	omiClient *Client
}

func NewMQClient(redisClient *redis.Client, namespace string, serverType string, prefix string) *MQClient {
	return &MQClient{omiClient: NewClient(redisClient, namespace, serverType, prefix)}
}

func (c *MQClient) NewConsumer(channel string, address string) *Consumer {
	return &Consumer{
		omiClient: c.omiClient,
		channel:   channel,
		address:   address,
		register:  c.omiClient.NewRegister(channel, address),
	}
}

func (c *MQClient) NewProducer(channel string) *Producer {
	producer := Producer{
		omiClient: c.omiClient,
		channel:   channel,
	}
	go producer.listen()
	return &producer
}

func (c *MQClient) GetOmiClient() *Client {
	return c.omiClient
}
