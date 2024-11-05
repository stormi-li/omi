package omi

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/go-redis/redis/v8"
)

func mapToJsonStr(data map[string]string) string {
	jsonStr, _ := json.MarshalIndent(data, " ", "  ")
	return string(jsonStr)
}

func jsonStrToMap(jsonStr string) map[string]string {
	var dataMap map[string]string
	json.Unmarshal([]byte(jsonStr), &dataMap)
	return dataMap
}

func splitCommand(address string) (string, string) {
	index := strings.Index(address, NamespaceSeparator)
	if index == -1 {
		return "", ""
	}
	return address[:index], address[index+1:]
}

func getKeysByNamespace(redisClient *redis.Client, namespace string) []string {
	var keys []string
	cursor := uint64(0)
	for {
		// 使用 SCAN 命令获取键名
		res, newCursor, err := redisClient.Scan(context.Background(), cursor, namespace+"*", 0).Result()
		if err != nil {
			return nil
		}
		// 处理键名，去掉命名空间
		for _, key := range res {
			// 去掉命名空间部分
			keyWithoutNamespace := key[len(namespace):]
			keys = append(keys, keyWithoutNamespace[1:])
		}
		cursor = newCursor
		// 如果游标为0，则结束循环
		if cursor == 0 {
			break
		}
	}
	return keys
}
