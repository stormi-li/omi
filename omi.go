package omi

import (
	"github.com/go-redis/redis/v8"
	ominager "github.com/stormi-li/omi/omi-manager"
	ominitor "github.com/stormi-li/omi/omi-monitor"
	omiweb "github.com/stormi-li/omi/omi-web"
)

func NewServerManager(opts *redis.Options) *ominager.Client {
	return ominager.NewClient(redis.NewClient(opts), ominager.Server, ominager.Prefix_Server)
}

func NewWebManager(opts *redis.Options) *ominager.Client {
	return ominager.NewClient(redis.NewClient(opts), ominager.Web, ominager.Prefix_Web)
}

func NewConfigManager(opts *redis.Options) *ominager.Client {
	return ominager.NewClient(redis.NewClient(opts), ominager.Config, ominager.Prefix_Config)
}

func NewWebClient(opts *redis.Options) *omiweb.Client {
	return omiweb.NewClient(redis.NewClient(opts), NewWebManager(opts), NewServerManager(opts))
}

func NewMonitor(opts *redis.Options) *ominitor.Client {
	return ominitor.NewClient(NewServerManager(opts).NewSearcher(), NewWebManager(opts).NewSearcher(), NewConfigManager(opts).NewSearcher())
}
