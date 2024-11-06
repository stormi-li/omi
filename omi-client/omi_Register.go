package omiclient

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	omipc "github.com/stormi-li/omi/omi-ipc"
)

type Register struct {
	redisClient      *redis.Client
	omipcClient      *omipc.Client
	serverName       string
	nodeType         string
	address          string
	redisChannelName string
	namespace        string
	ctx              context.Context
}

func (register *Register) StartOnMain(data ...map[string]string) {
	if len(data) == 0 {
		go register.start(node_main, map[string]string{})
	}
	go register.start(node_main, data[0])
}

func (register *Register) StartOnBackup(data ...map[string]string) {
	if len(data) == 0 {
		register.start(node_backup, map[string]string{})
	}
	register.start(node_backup, data[0])
}

func (register *Register) start(nodeType string, data map[string]string) {
	jsonStrData := mapToJsonStr(data)
	register.nodeType = nodeType
	nodeState := state_start

	go func() {
		for {
			key := register.namespace + register.serverName + omipc.NamespaceSeparator + nodeState + omipc.NamespaceSeparator + nodeType + omipc.NamespaceSeparator + register.address
			register.redisClient.Set(register.ctx, key, jsonStrData, const_expireTime)
			time.Sleep(const_expireTime / 2)
		}
	}()
	log.Println("register server for", register.serverName+"["+register.address+"]", "is starting")
	channel := register.serverName + omipc.NamespaceSeparator + register.address
	listener := register.omipcClient.NewListener(channel)
	listener.Listen(func(msg string) {
		if msg == command_start {
			nodeState = state_start
		}
		if msg == command_stop {
			nodeState = state_stop
		}
		if msg == command_toMain {
			nodeType = node_main
		}
		if msg == command_toBackup {
			nodeType = node_backup
		}
	})
}

func (register *Register) ToMain() {
	register.omipcClient.Notify(register.redisChannelName, command_toMain)
}

func (register *Register) ToBackup() {
	register.omipcClient.Notify(register.redisChannelName, command_toBackup)
}

func (register *Register) ToStart() {
	register.omipcClient.Notify(register.redisChannelName, command_start)
}

func (register *Register) ToStop() {
	register.omipcClient.Notify(register.redisChannelName, command_stop)
}
