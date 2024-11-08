package ominager

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi"
	omiclient "github.com/stormi-li/omi/omi-client"
)

type Manager struct {
	serverClient   *omiclient.Client
	webClient      *omiclient.Client
	configClient   *omiclient.Client
	serverSearcher *omiclient.Searcher
	webSearcher    *omiclient.Searcher
	configSearcher *omiclient.Searcher
}

func NewManager(opts *redis.Options) *Manager {
	serverClient := omi.NewServerClient(opts)
	webClient := omi.NewWebClient(opts)
	configClient := omi.NewConfigClient(opts)
	return &Manager{
		serverClient:   serverClient,
		webClient:      webClient,
		configClient:   configClient,
		serverSearcher: serverClient.NewSearcher(),
		webSearcher:    webClient.NewSearcher(),
		configSearcher: configClient.NewSearcher(),
	}
}

func (manager *Manager) GetServerNodes() map[string]map[string]map[string]string {
	return manager.serverSearcher.AllServers()
}

func (manager *Manager) GetWebNodes() map[string]map[string]map[string]string {
	return manager.webSearcher.AllServers()
}

func (manager *Manager) GetConfigNodes() map[string]map[string]map[string]string {
	return manager.configSearcher.AllServers()
}

func (manager *Manager) Handler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	// 获取请求的路径并去掉开头的 '/'
	path := strings.TrimPrefix(r.URL.Path, "/")
	// 以 '/' 分割路径
	parts := strings.Split(path, "/")

	if parts[0] == command_GetWebNodes {
		w.Write([]byte(mapToJsonStr(manager.GetWebNodes())))
	}
	if parts[0] == command_GetServerNodes {
		w.Write([]byte(mapToJsonStr(manager.GetServerNodes())))
	}
	if parts[0] == command_GetConfigNodes {
		w.Write([]byte(mapToJsonStr(manager.GetConfigNodes())))
	}
}

func mapToJsonStr(data map[string]map[string]map[string]string) string {
	jsonStr, _ := json.MarshalIndent(data, " ", "  ")
	return string(jsonStr)
}
