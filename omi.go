package omi

import (
	"github.com/go-redis/redis/v8"
	omiclient "github.com/stormi-li/omi/omi-client"
	ominager "github.com/stormi-li/omi/omi-manager"
	omiweb "github.com/stormi-li/omi/omi-web"
)

func NewServerClient(opts *redis.Options) *omiclient.Client {
	return omiclient.NewClient(redis.NewClient(opts), omiclient.Server, omiclient.Prefix_Server)
}

func NewWebClient(opts *redis.Options) *omiclient.Client {
	return omiclient.NewClient(redis.NewClient(opts), omiclient.Web, omiclient.Prefix_Web)
}

func NewConfigClient(opts *redis.Options) *omiclient.Client {
	return omiclient.NewClient(redis.NewClient(opts), omiclient.Config, omiclient.Prefix_Config)
}

func NewOmiweb(opts *redis.Options) *omiweb.Client {
	return omiweb.NewClient(redis.NewClient(opts), NewWebClient(opts))
}

func NewManager(opts *redis.Options) *ominager.Client {
	return ominager.NewClient(NewServerClient(opts).NewSearcher(), NewWebClient(opts).NewSearcher(), NewConfigClient(opts).NewSearcher())
}
