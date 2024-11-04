package omi

import (
	"context"
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	omipc "github.com/stormi-li/omi/om-ipc"
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
	sort.Slice(addrs, func(a, b int) bool {
		return addrs[a] < addrs[b]
	})
	return addrs
}

func (searcher *Searcher) AllServers() []string {
	addrs := getKeysByNamespace(searcher.redisClient, searcher.namespace[:len(searcher.namespace)-1])
	sort.Slice(addrs, func(a, b int) bool {
		return addrs[a] < addrs[b]
	})
	return addrs
}

func (searcher *Searcher) GetHighestPriorityServer(serverName string) (string, map[string]string) {
	addrs := searcher.SearchStartingServers(serverName)
	var validAddr string
	if len(addrs) > 0 {
		validAddr = split(addrs[0])[1]
		data, _ := searcher.redisClient.Get(searcher.ctx, searcher.namespace+serverName+const_separator+addrs[0]).Result()
		json.Unmarshal([]byte(data), &searcher.data)
	}
	return validAddr, searcher.data
}

func (searcher *Searcher) GetData(serverName, state, nodeType, address string) map[string]string {
	key := searcher.namespace + serverName + const_separator + state + const_separator + nodeType + const_separator + address
	data, _ := searcher.redisClient.Get(searcher.ctx, key).Result()
	var dataMap map[string]string
	json.Unmarshal([]byte(data), &dataMap)
	return dataMap
}

func (searcher *Searcher) Listen(serverName string, handler func(address string, data map[string]string)) {
	addr := ""
	jsonByte, _ := json.MarshalIndent(searcher.data, " ", "  ")
	dataStr := string(jsonByte)
	for {
		newAddr, data := searcher.GetHighestPriorityServer(serverName)
		jsonByte, _ = json.MarshalIndent(data, " ", "  ")
		newDataStr := string(jsonByte)
		if newAddr != addr || newDataStr != dataStr {
			addr = newAddr
			dataStr = newDataStr
			handler(addr, searcher.data)
		}
		time.Sleep(2 * time.Second)
	}
}

func (searcher *Searcher) SearchStartingServers(serverName string) []string {
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
	index := strings.Index(address, const_separator)
	if index == -1 {
		return nil
	}
	return []string{address[:index], address[index+1:]}
}