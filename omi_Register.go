package omi

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	omipc "github.com/stormi-li/omi/omi-ipc"
)

type Register struct {
	redisClient      *redis.Client
	omipcClient      *omipc.Client
	namespace        string
	ctx              context.Context
	serverName       string
	nodeType         string
	address          string
	redisChannelName string
}

func (register *Register) StartOnMain(data map[string]string) {
	register.start(node_main, data)
}

func (register *Register) StartOnBackup(data map[string]string) {
	register.start(node_backup, data)
}

func (register *Register) start(nodeType string, data map[string]string) {
	jsonByte, _ := json.MarshalIndent(data, " ", "  ")
	jsonStrData := string(jsonByte)
	register.nodeType = nodeType
	nodeState := state_start

	go func() {
		for {
			key := register.namespace + register.serverName + NamespaceSeparator + nodeState + NamespaceSeparator + nodeType + NamespaceSeparator + register.address
			register.redisClient.Set(register.ctx, key, jsonStrData, const_expireTime)
			time.Sleep(const_expireTime / 2)
		}
	}()
	log.Println("register server for", register.serverName+"["+register.address+"]", "is starting")
	channel := register.serverName + NamespaceSeparator + register.address
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
		if command, json := splitCommand(msg); command == command_updateNodeData {
			jsonStrData = json
		}
	})
}

func (register *Register) ToMain() {
	register.omipcClient.Notify(register.redisChannelName, command_toMain)
}

func (register *Register) ToBackup() {
	register.omipcClient.Notify(register.redisChannelName, command_toBackup)
}

func (register *Register) Start() {
	register.omipcClient.Notify(register.redisChannelName, command_start)
}

func (register *Register) Stop() {
	register.omipcClient.Notify(register.redisChannelName, command_stop)
}

func (register *Register) UpdateData(data map[string]string) {
	register.omipcClient.Notify(register.redisChannelName, command_updateNodeData+":"+mapToJsonStr(data))
}
