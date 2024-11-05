package omiclient

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	omipc "github.com/stormi-li/omi/omi-ipc"
)

type Searcher struct {
	redisClient *redis.Client
	omipcClient *omipc.Client
	namespace   string
	ctx         context.Context
	data        map[string]string
}

func (searcher *Searcher) SearchAllServers(serverName string) []string {
	addrs := getKeysByNamespace(searcher.redisClient, searcher.namespace+serverName)
	sort.Slice(addrs, func(i, j int) bool {
		return addrs[i] > addrs[j]
	})
	return addrs
}

func (searcher *Searcher) AllServers() []string {
	return getKeysByNamespace(searcher.redisClient, searcher.namespace[:len(searcher.namespace)-1])
}

func (searcher *Searcher) GetHighestPriorityServer(serverName string) (string, map[string]string) {
	addrs := searcher.SearchRunningServers(serverName)
	var validAddr string
	if len(addrs) > 0 {
		validAddr = split(addrs[0])[1]
		data, _ := searcher.redisClient.Get(searcher.ctx, searcher.namespace+serverName+omipc.NamespaceSeparator+addrs[0]).Result()
		searcher.data = jsonStrToMap(data)
	}
	return validAddr, searcher.data
}

func (searcher *Searcher) GetData(serverName, state, nodeType, address string) map[string]string {
	key := searcher.namespace + serverName + omipc.NamespaceSeparator + state + omipc.NamespaceSeparator + nodeType + omipc.NamespaceSeparator + address
	data, _ := searcher.redisClient.Get(searcher.ctx, key).Result()
	return jsonStrToMap(data)
}

func (searcher *Searcher) Listen(serverName string, handler func(address string, data map[string]string)) {
	addr := ""
	for {
		newAddr, _ := searcher.GetHighestPriorityServer(serverName)
		if newAddr != addr {
			addr = newAddr
			handler(addr, searcher.data)
		}
		time.Sleep(const_listenWaitTime)
	}
}

func (searcher *Searcher) SearchRunningServers(serverName string) []string {
	servers := searcher.SearchAllServers(serverName)
	startingservers := []string{}
	for _, val := range servers {
		temp := split(val)
		if temp[0] == state_start {
			startingservers = append(startingservers, temp[1])
		}
	}
	return startingservers
}

func split(address string) []string {
	index := strings.Index(address, omipc.NamespaceSeparator)
	if index == -1 {
		return nil
	}
	return []string{address[:index], address[index+1:]}
}
