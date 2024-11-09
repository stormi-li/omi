package ominager

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	omiclient "github.com/stormi-li/omi/omi-client"
)

type Manager struct {
	serverSearcher *omiclient.Searcher
	webSearcher    *omiclient.Searcher
	configSearcher *omiclient.Searcher
}

func NewManager(serverSearcher *omiclient.Searcher, webSearcher *omiclient.Searcher, configSearcher *omiclient.Searcher) *Manager {
	return &Manager{
		serverSearcher: serverSearcher,
		webSearcher:    webSearcher,
		configSearcher: configSearcher,
	}
}

func (manager *Manager) GetServerNodes() map[string]map[string]map[string]string {
	return manager.serverSearcher.SearchAllServers()
}

func (manager *Manager) GetWebNodes() map[string]map[string]map[string]string {
	return manager.webSearcher.SearchAllServers()
}

func (manager *Manager) GetConfigNodes() map[string]map[string]map[string]string {
	return manager.configSearcher.SearchAllServers()
}

func (manager *Manager) Handler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	// 获取请求的路径并去掉开头的 '/'
	path := strings.TrimPrefix(r.URL.Path, "/")
	// 以 '/' 分割路径
	parts := strings.Split(path, "/")

	if parts[0] == command_GetWebNodes {
		w.Write([]byte(toJsonStr(manager.GetWebNodes())))
	}
	if parts[0] == command_GetServerNodes {
		w.Write([]byte(toJsonStr(manager.GetServerNodes())))
	}
	if parts[0] == command_GetConfigNodes {
		w.Write([]byte(toJsonStr(manager.GetConfigNodes())))
	}
}

func toJsonStr(nodes map[string]map[string]map[string]string) string {
	res := [][]string{}
	for name, addresses := range nodes {
		for address, details := range addresses {
			weight := details["weight"]
			res = append(res, []string{name, address, weight})
		}
	}
	return sliceToJsonStr(res)
}

func sliceToJsonStr(data [][]string) string {
	jsonStr, _ := json.MarshalIndent(data, " ", "  ")
	return string(jsonStr)
}
