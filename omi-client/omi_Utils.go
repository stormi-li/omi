package omiclient

import (
	"context"
	"encoding/json"

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

func getKeysByNamespace(redisClient *redis.Client, namespace string) []string {
	var keys []string
	cursor := uint64(0)
	for {
		res, newCursor, err := redisClient.Scan(context.Background(), cursor, namespace+"*", 0).Result()
		if err != nil {
			return nil
		}
		for _, key := range res {
			keyWithoutNamespace := key[len(namespace):]
			keys = append(keys, keyWithoutNamespace[1:])
		}
		cursor = newCursor
		if cursor == 0 {
			break
		}
	}
	return keys
}
