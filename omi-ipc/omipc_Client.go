package omipc

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Client struct {
	redisClient *redis.Client
	ctx         context.Context
}

func NewOmipc(redisClient *redis.Client) *Client {
	return &Client{
		redisClient: redisClient,
		ctx:         context.Background(),
	}
}

func (c *Client) Notify(channel, msg string) {
	c.redisClient.Publish(c.ctx, channel, msg)
}

func (c *Client) Wait(channel string, timeout time.Duration) string {
	sub := c.redisClient.Subscribe(c.ctx, channel)

	msgChan := sub.Channel()
	defer sub.Close()

	if timeout == 0 {
		msg := <-msgChan
		return msg.Payload
	}

	timer := time.NewTicker(timeout)
	defer timer.Stop()

	select {
	case <-timer.C:
		return ""
	case msg := <-msgChan:
		return msg.Payload
	}
}

func (c *Client) NewListener(channel string) *Listener {
	sub := c.redisClient.Subscribe(c.ctx, channel)
	return &Listener{
		shutdown: make(chan struct{}, 1),
		sub:      sub,
	}
}

func (c *Client) NewLock(lockName string) *Lock {
	return &Lock{
		uuid:        uuid.NewString(),
		lockName:    lockName,
		stop:        make(chan struct{}, 1),
		omipcClient: c,
		redisClient: c.redisClient,
		ctx:         context.Background(),
	}
}
