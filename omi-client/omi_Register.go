package omiclient

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	omipc "github.com/stormi-li/omi/omi-ipc"
)

type Register struct {
	redisClient      *redis.Client
	omipcClient      *omipc.Client
	serverName       string
	address          string
	redisChannelName string
	namespace        string
	ctx              context.Context
}

func (register *Register) Start(weight int, data map[string]string) {
	data["weight"] = strconv.Itoa(weight)
	jsonStrData := mapToJsonStr(data)
	go func() {
		for {
			key := register.namespace + register.serverName + namespace_separator + register.address
			register.redisClient.Set(register.ctx, key, jsonStrData, const_expireTime)
			time.Sleep(const_expireTime / 2)
		}
	}()
	log.Println("register server for", register.serverName+"["+register.address+"]", "is starting")
}
