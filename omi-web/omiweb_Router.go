package omiweb

import (
	"math/rand/v2"
	"strconv"
	"sync"
	"time"

	omiclient "github.com/stormi-li/omi/omi-client"
)

type Router struct {
	searcher   *omiclient.Searcher
	addressMap map[string][]string
	mutex      sync.Mutex
}

func (router *Router) refresh() {
	nodeMap := router.searcher.SearchAllServers()
	addrMap := map[string][]string{}
	for name, addrs := range nodeMap {
		for addr, data := range addrs {
			weight, _ := strconv.Atoi(data["weight"])
			for i := 0; i < weight; i++ {
				addrMap[name] = append(addrMap[name], addr)
			}
		}
	}
	router.mutex.Lock()
	router.addressMap = addrMap
	router.mutex.Unlock()
}

func newRouter(searcher *omiclient.Searcher) *Router {
	router := Router{
		searcher:   searcher,
		addressMap: map[string][]string{},
		mutex:      sync.Mutex{},
	}
	go func() {
		for {
			router.refresh()
			time.Sleep(router_refresh_interval)
		}
	}()
	return &router
}

func (router *Router) getAddress(serverName string) string {
	router.mutex.Lock()
	defer router.mutex.Unlock()
	if len(router.addressMap[serverName]) == 0 {
		return ""
	}
	return router.addressMap[serverName][rand.IntN(len(router.addressMap[serverName]))]
}

func (router *Router) Has(serverName string) bool {
	return len(router.addressMap[serverName]) != 0
}
