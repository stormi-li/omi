package ominager

import (
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi"
	omiclient "github.com/stormi-li/omi/omi-client"
)

type Manager struct {
	serverClient   *omiclient.Client
	mqClient       *omiclient.MQClient
	configClient   *omiclient.Client
	serverSearcher *omiclient.Searcher
	mqSearcher     *omiclient.Searcher
	configSearcher *omiclient.Searcher
	nodeMap        map[string]Node
}

func NewManager(redisClient *redis.Client, namespace string) *Manager {
	serverClient := omi.NewServerClient(redisClient, namespace)
	mqClient := omi.NewMQClient(redisClient, namespace)
	configClient := omi.NewConfigClient(redisClient, namespace)
	return &Manager{
		serverClient:   serverClient,
		mqClient:       mqClient,
		configClient:   configClient,
		serverSearcher: serverClient.NewSearcher(),
		mqSearcher:     mqClient.GetOmiClient().NewSearcher(),
		configSearcher: configClient.NewSearcher(),
		nodeMap:        map[string]Node{},
	}
}

func (manager *Manager) GetServerNodes() []Node {
	return manager.toNodeSlice(server, manager.serverClient, manager.serverSearcher)
}

func (manager *Manager) GetMQNodes() []Node {
	return manager.toNodeSlice(mq, manager.mqClient.GetOmiClient(), manager.mqSearcher)
}

func (manager *Manager) GetConfigNodes() []Node {
	return manager.toNodeSlice(config, manager.configClient, manager.configSearcher)
}

func (manager *Manager) toNodeSlice(serverType string, omiClient *omiclient.Client, searcher *omiclient.Searcher) []Node {
	keys := searcher.AllServers()
	nodes := []Node{}

	for _, val := range keys {
		info := spliteNodeKey(val)
		node := *newNode(serverType, info[0], info[1], info[2], info[3], omiClient, searcher)
		nodes = append(nodes, node)
		manager.nodeMap[info[0]+":"+info[3]] = node
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Address > nodes[j].Address
	})
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].ServerName < nodes[j].ServerName
	})
	return nodes
}

func (manager *Manager) Handler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	// 获取请求的路径并去掉开头的 '/'
	path := strings.TrimPrefix(r.URL.Path, "/")
	// 以 '/' 分割路径
	parts := strings.Split(path, "/")

	parts = parts[1:]
	if parts[0] == command_GetMQNodes {
		w.Write([]byte(nodesToString(manager.GetMQNodes())))
	}
	if parts[0] == command_GetServerNodes {
		w.Write([]byte(nodesToString(manager.GetServerNodes())))
	}
	if parts[0] == command_GetConfigNodes {
		w.Write([]byte(nodesToString(manager.GetConfigNodes())))
	}
	if parts[0] == command_GetAllNodes {
		nodes := manager.GetServerNodes()
		nodes = append(nodes, manager.GetMQNodes()...)
		nodes = append(nodes, manager.GetConfigNodes()...)
		w.Write([]byte(nodesToString(nodes)))
	}

	getNode := func() *Node {
		key := parts[1] + ":" + parts[2]
		node := manager.nodeMap[key]
		if node.Address == "" {
			manager.GetServerNodes()
			manager.GetMQNodes()
			manager.GetConfigNodes()
		}
		node = manager.nodeMap[key]
		return &node
	}

	if parts[0] == command_ToMain {
		node := getNode()
		node.ToMain()
		w.Write([]byte(node.ToString()))
	}
	if parts[0] == command_ToBackup {
		node := getNode()
		node.ToBackup()
		w.Write([]byte(node.ToString()))
	}
	if parts[0] == command_Stop {
		node := getNode()
		node.Stop()
		w.Write([]byte(node.ToString()))
	}
	if parts[0] == command_Start {
		node := getNode()
		node.Start()
		w.Write([]byte(node.ToString()))
	}
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
	index := strings.Index(address, ":")
	if index == -1 {
		return nil
	}
	return []string{address[:index], address[index+1:]}
}
