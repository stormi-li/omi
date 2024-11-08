package omiweb

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/stormi-li/omi"
)

func NewClient(redisClient *redis.Client, namespace string) *Client {
	omiClient := omi.NewServerClient(redisClient, namespace)
	return &Client{
		router:             newRouter(omiClient.NewSearcher()),
		redisClient:        redisClient,
		omiClient:          omiClient,
		namespace:          namespace,
		originCheckHandler: func(r *http.Request) bool { return true },
		upgrader:           websocket.Upgrader{},
	}
}
