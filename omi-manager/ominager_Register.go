package omiclient

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type Register struct {
	redisClient *redis.Client
	serverName  string
	address     string
	weight      int
	Data        map[string]string
	namespace   string
	ctx         context.Context
}

func (register *Register) Register(address string) {
	register.address = address
	register.Data["weight"] = strconv.Itoa(register.weight)
	jsonStrData := mapToJsonStr(register.Data)
	go func() {
		for {
			key := register.namespace + register.serverName + namespace_separator + register.address
			register.redisClient.Set(register.ctx, key, jsonStrData, const_expireTime)
			time.Sleep(const_expireTime / 2)
		}
	}()
	log.Println("register server for", register.serverName+"["+register.address+"]", "is starting")
}
