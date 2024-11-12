package omiclient

import (
	"context"
	"strconv"
	"strings"

	"math/rand"

	"github.com/go-redis/redis/v8"
)

type Searcher struct {
	redisClient *redis.Client
	namespace   string
	ctx         context.Context
}

func (searcher *Searcher) SearchByName(serverName string) map[string]map[string]string {
	keys := getKeysByNamespace(searcher.redisClient, searcher.namespace+serverName)
	res := map[string]map[string]string{}
	for _, key := range keys {
		data, _ := searcher.redisClient.Get(searcher.ctx, searcher.namespace+serverName+namespace_separator+key).Result()
		res[key] = jsonStrToMap(data)
	}
	return res
}

func (searcher *Searcher) IsAlive(serverName string, address string) bool {
	data, _ := searcher.redisClient.Get(searcher.ctx, searcher.namespace+serverName+namespace_separator+address).Result()
	return data != ""
}

func (searcher *Searcher) SearchAllServers() map[string]map[string]map[string]string {
	keys := getKeysByNamespace(searcher.redisClient, searcher.namespace[:len(searcher.namespace)-1])
	res := map[string]map[string]map[string]string{}
	for _, key := range keys {
		data, _ := searcher.redisClient.Get(searcher.ctx, searcher.namespace+key).Result()
		parts := split(key)
		if res[parts[0]] == nil {
			res[parts[0]] = map[string]map[string]string{}
		}
		res[parts[0]][parts[1]] = jsonStrToMap(data)
	}
	return res
}

func (searcher *Searcher) SearchByLoadBalancing(serverName string) (string, map[string]string) {
	addrs := searcher.SearchByName(serverName)
	var addressPool []string
	var dataPool []map[string]string
	for name, data := range addrs {
		weight, _ := strconv.Atoi(data["weight"])
		for i := 0; i < weight; i++ {
			addressPool = append(addressPool, name)
			dataPool = append(dataPool, data)
		}
	}
	if len(addressPool) == 0 {
		return "", nil
	}
	selectIndex := rand.Intn(len(addressPool))
	return addressPool[selectIndex], dataPool[selectIndex]
}

func split(address string) []string {
	index := strings.Index(address, namespace_separator)
	if index == -1 {
		return nil
	}
	return []string{address[:index], address[index+1:]}
}
