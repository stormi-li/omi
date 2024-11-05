package ominager

import (
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi"
)

type Manager struct {
	serverSearcher *omi.Searcher
	mqSearcher     *omi.Searcher
	configSearcher *omi.Searcher
	serverClient   *omi.Client
	mqClient       *omi.Client
	configClient   *omi.Client
	nodeMap        map[string]Node
}

func NewManager(redisClient *redis.Client, namespace string) *Manager {
	serverClient := omi.NewClient(redisClient, namespace, omi.Server)
	mqClient := omi.NewClient(redisClient, namespace, omi.MQ)
	configClient := omi.NewClient(redisClient, namespace, omi.Config)
	return &Manager{
		serverClient:   serverClient,
		mqClient:       mqClient,
		configClient:   configClient,
		serverSearcher: serverClient.NewSearcher(),
		mqSearcher:     mqClient.NewSearcher(),
		configSearcher: configClient.NewSearcher(),
		nodeMap:        map[string]Node{},
	}
}

func (manager *Manager) GetServerNodes() []Node {
	return manager.toNodeSlice(omi.Server)
}

func (manager *Manager) GetMQNodes() []Node {
	return manager.toNodeSlice(omi.MQ)
}

func (manager *Manager) GetConfigNodes() []Node {
	return manager.toNodeSlice(omi.Config)
}
func (manager *Manager) toNodeSlice(serverType omi.ServerType) []Node {
	var keys []string
	nodes := []Node{}
	var client *omi.Client
	var searcher *omi.Searcher
	if serverType == omi.Server {
		client = manager.serverClient
		searcher = manager.serverSearcher
		keys = manager.serverSearcher.AllServers()
	}
	if serverType == omi.MQ {
		client = manager.mqClient
		searcher = manager.mqSearcher
		keys = manager.mqSearcher.AllServers()
	}
	if serverType == omi.Config {
		client = manager.configClient
		searcher = manager.configSearcher
		keys = manager.configSearcher.AllServers()
	}
	for _, val := range keys {
		info := spliteNodeKey(val)
		node := *newNode(serverType, info[0], info[1], info[2], info[3], client, searcher)
		nodes = append(nodes, node)
		manager.nodeMap[info[0]+const_separator+info[3]] = node
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].ServerName < nodes[j].ServerName
	})
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Address < nodes[j].Address
	})
	return nodes
}

func spliteNodeKey(key string) []string {
	res := []string{}
	for i := 0; i < 3; i++ {
		temp := split(key)
		key = temp[1]
		res = append(res, temp[0])
	}
	res = append(res, key)
	return res
}

func split(address string) []string {
	index := strings.Index(address, const_separator)
	if index == -1 {
		return nil
	}
	return []string{address[:index], address[index+1:]}
}

func (manager *Manager) Handler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	// 获取请求的路径并去掉开头的 '/'
	path := strings.TrimPrefix(r.URL.Path, "/")
	// 以 '/' 分割路径
	parts := strings.Split(path, "/")

	parts = parts[1:]
	if parts[0] == "GetMQNodes" {
		w.Write([]byte(nodesToString(manager.GetMQNodes())))
	}
	if parts[0] == "GetServerNodes" {
		w.Write([]byte(nodesToString(manager.GetServerNodes())))
	}
	if parts[0] == "GetConfigNodes" {
		w.Write([]byte(nodesToString(manager.GetConfigNodes())))
	}
	if parts[0] == "GetAllNodes" {
		nodes := manager.GetServerNodes()
		nodes = append(nodes, manager.GetMQNodes()...)
		nodes = append(nodes, manager.GetConfigNodes()...)
		w.Write([]byte(nodesToString(nodes)))
	}

	getNode := func() *Node {
		key := parts[1] + const_separator + parts[2]
		node := manager.nodeMap[key]
		if node.Address == "" {
			manager.GetServerNodes()
			manager.GetMQNodes()
			manager.GetConfigNodes()
		}
		node = manager.nodeMap[key]
		return &node
	}

	if parts[0] == "ToMain" {
		node := getNode()
		node.ToMain()
		w.Write([]byte(node.ToString()))
	}
	if parts[0] == "ToBackup" {
		node := getNode()
		node.ToBackup()
		w.Write([]byte(node.ToString()))
	}
	if parts[0] == "Stop" {
		node := getNode()
		node.Stop()
		w.Write([]byte(node.ToString()))
	}
	if parts[0] == "Start" {
		node := getNode()
		node.Start()
		w.Write([]byte(node.ToString()))
	}
}
