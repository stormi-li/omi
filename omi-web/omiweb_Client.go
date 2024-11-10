package omiweb

import (
	"github.com/go-redis/redis/v8"
	omiclient "github.com/stormi-li/omi/omi-client"
)

type Client struct {
	redisClient     *redis.Client
	omiWebClient    *omiclient.Client
	omiServerClient *omiclient.Client
}

func (c *Client) NewWebServer(serverName string, weight int) *WebServer {
	return newWebServer(c.redisClient, c.omiWebClient, c.omiServerClient, serverName, weight)
}

func (c *Client) GenerateTemplate() {
	copyResource(getSourceFilePath() + source_path)
}

func (c *Client) NewProxyServer(serverName string) *ProxyServer {
	return &ProxyServer{
		router:       newRouter(c.omiWebClient.NewSearcher()),
		omiWebClient: c.omiWebClient,
		serverName:   serverName,
	}
}

func (c *Client) DisableLog() {
	log_cache = false
}
