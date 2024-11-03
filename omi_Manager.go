package omi

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
)

type Manager struct {
	serverSearcher      *Searcher
	mqSearcher          *Searcher
	configSearcher      *Searcher
	serverReseardClient *Client
	mqReseardClient     *Client
	configReseardClient *Client
	nodeMap             map[string]Node
}

func newManager(redisClient *redis.Client, namespace string) *Manager {
	serverReseardClient := NewClient(redisClient, namespace, Server)
	mqReseardClient := NewClient(redisClient, namespace, MQ)
	configReseardClient := NewClient(redisClient, namespace, Config)
	return &Manager{
		serverReseardClient: serverReseardClient,
		mqReseardClient:     mqReseardClient,
		configReseardClient: configReseardClient,
		serverSearcher:      serverReseardClient.NewSearcher(),
		mqSearcher:          mqReseardClient.NewSearcher(),
		configSearcher:      configReseardClient.NewSearcher(),
		nodeMap:             map[string]Node{},
	}
}

func (manager *Manager) GetServerNodes() []Node {
	return manager.toNodeSlice(Server)
}

func (manager *Manager) GetMQNodes() []Node {
	return manager.toNodeSlice(MQ)
}

func (manager *Manager) GetConfigNodes() []Node {
	return manager.toNodeSlice(Config)
}
func (manager *Manager) toNodeSlice(serverType ServerType) []Node {
	var keys []string
	nodes := []Node{}
	var client *Client
	var searcher *Searcher
	if serverType == Server {
		client = manager.serverReseardClient
		searcher = manager.serverSearcher
		keys = manager.serverSearcher.allServers()
	}
	if serverType == MQ {
		client = manager.mqReseardClient
		searcher = manager.mqSearcher
		keys = manager.mqSearcher.allServers()
	}
	if serverType == Config {
		client = manager.configReseardClient
		searcher = manager.configSearcher
		keys = manager.configSearcher.allServers()
	}
	for _, val := range keys {
		info := spliteNodeKey(val)
		node := *newNode(serverType, info[0], info[1], info[2], info[3], client, searcher)
		nodes = append(nodes, node)
		manager.nodeMap[info[0]+const_separator+info[3]] = node
	}
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

func (manager *Manager) handler(w http.ResponseWriter, r *http.Request) {
	// 获取请求的路径并去掉开头的 '/'
	path := strings.TrimPrefix(r.URL.Path, "/")

	// 以 '/' 分割路径，获取第一个参数
	parts := strings.Split(path, "/")
	if parts[0] == "GetMQNodes" {
		w.Write([]byte(nodesToString(manager.GetMQNodes())))
	}
	if parts[0] == "GetServerNodes" {
		w.Write([]byte(nodesToString(manager.GetServerNodes())))
	}
	if parts[0] == "GetConfigNodes" {
		w.Write([]byte(nodesToString(manager.GetConfigNodes())))
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

	if parts[0] == "GetNodeData" {
		node := getNode()
		_, data := node.GetData()
		w.Write([]byte(data))
	}
	if parts[0] == "ToMain" {
		node := getNode()
		node.ToMain()
		w.Write([]byte(node.ToString()))
	}
	if parts[0] == "ToStandby" {
		node := getNode()
		node.ToStandby()
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
	if parts[0] == "Close" {
		node := getNode()
		node.Close()
		w.Write([]byte(node.ToString()))
	}
}

func (manager *Manager) Start(managerName, address string) {
	register := manager.serverReseardClient.NewRegister(managerName, address)
	go register.StartOnMain(map[string]string{"message": "omi manager server"})
	http.HandleFunc("/", manager.handler)
	log.Println("omi manager server: " + managerName + " is running on http://" + address)
	go http.ListenAndServe(":"+strings.Split(address, ":")[1], nil)
	<-register.CloseSignal
}
