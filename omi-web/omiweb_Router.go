package omiweb

import (
	omiclient "github.com/stormi-li/omi/omi-client"
)

type Router struct {
	searcher   *omiclient.Searcher
	addressMap map[string][]string
}

func newRouter(searcher *omiclient.Searcher) *Router {
	return &Router{
		searcher:   searcher,
		addressMap: map[string][]string{},
	}
}

func (router *Router) getAddress(serverName string) string {
	if len(router.addressMap[serverName]) != 2 {
		address, _ := router.searcher.SearchOneByWeight(serverName)
		if address != "" {
			router.addressMap[serverName] = []string{address, getCurrentTimeString()}
		} else {
			return ""
		}
	}
	go router.refresh(serverName)
	return router.addressMap[serverName][0]
}

func (router *Router) refresh(serverName string) {
	if isMoreThanTwoSecondsAgo(router.addressMap[serverName][1]) {
		address, _ := router.searcher.SearchOneByWeight(serverName)
		router.addressMap[serverName][0] = address
		router.addressMap[serverName][1] = getCurrentTimeString()
	}
}
