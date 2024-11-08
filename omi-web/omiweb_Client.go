package omiweb

import (
	"github.com/go-redis/redis/v8"
	omiclient "github.com/stormi-li/omi/omi-client"
)

type Client struct {
	redisClient *redis.Client
	omiClient   *omiclient.Client
}

func (omiweb *Client) NewWebServer(serverName string) *WebServer {
	return newWebServer(omiweb.redisClient, serverName)
}

func (omiweb *Client) GenerateTemplate() {
	copyResource(getSourceFilePath() + source_path)
}
