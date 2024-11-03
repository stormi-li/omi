package omiweb

import (
	"time"

	"github.com/stormi-li/omi"
)

type Router struct {
	searcher   *omi.Searcher
	addressMap map[string][]string
}

func newRouter(searcher *omi.Searcher) *Router {
	return &Router{
		searcher:   searcher,
		addressMap: map[string][]string{},
	}
}

func (router *Router) getAddress(serverName string) string {
	if len(router.addressMap[serverName]) != 2 {
		address, _ := router.searcher.GetHighestPriorityServer(serverName)
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
		address, _ := router.searcher.GetHighestPriorityServer(serverName)
		router.addressMap[serverName][0] = address
		router.addressMap[serverName][1] = getCurrentTimeString()
	}
}

// 获取当前时间并转换为 UTC 字符串
func getCurrentTimeString() string {
	currentTime := time.Now().UTC() // 设置为 UTC
	return currentTime.Format("2006-01-02 15:04:05")
}

// 将时间字符串解析为 UTC 时间
func parseTimeString(timeString string) (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	parsedTime, err := time.Parse(layout, timeString)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime.UTC(), nil // 设置为 UTC
}

// 比较字符串时间和当前时间，判断是否超过 2 秒
func isMoreThanTwoSecondsAgo(timeString string) bool {
	parsedTime, err := parseTimeString(timeString)
	if err != nil {
		return true // 如果解析出错，直接返回 true
	}

	currentTime := time.Now().UTC() // 统一设置为 UTC
	twoSecondsLater := parsedTime.Add(2 * time.Second)

	return currentTime.After(twoSecondsLater)
}
