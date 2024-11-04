package omi

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	omipc "github.com/stormi-li/omi/om-ipc"
)

type Register struct {
	redisClient *redis.Client
	omipcClient *omipc.Client
	namespace   string
	ctx         context.Context
	serverName  string
	nodeType    string
	address     string
	channel     string
	CloseSignal chan struct{}
}

func (register *Register) StartOnMain(data map[string]string) {
	register.start(node_main, data)
}

func (register *Register) StartOnBackup(data map[string]string) {
	register.start(node_backup, data)
}

func (register *Register) start(nodeType string, data map[string]string) {
	jsonByte, _ := json.MarshalIndent(data, " ", "  ")
	jsonStr := string(jsonByte)
	register.nodeType = nodeType
	nodeState := state_start
	ticker := time.NewTicker(const_expireTime / 2)
	close := make(chan struct{}, 1)
	update := make(chan struct{}, 1)

	updateHandler := func() {
		key := register.namespace + register.serverName + const_separator + nodeState + const_separator + nodeType + const_separator + register.address
		register.redisClient.Set(register.ctx, key, jsonStr, const_expireTime)
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				updateHandler()
			case <-update:
				updateHandler()
			case <-close:
				return
			}
		}
	}()
	log.Println("register server for", register.serverName+"["+register.address+"]", "is starting")
	channel := register.serverName + const_separator + register.address
	listener := register.omipcClient.NewListener(channel)
	listener.Listen(func(msg string) {
		if msg == command_close {
			close <- struct{}{}
			register.CloseSignal <- struct{}{}
			listener.Close()
		}
		if msg == command_start {
			nodeState = state_start
			update <- struct{}{}
		}
		if msg == command_stop {
			nodeState = state_stop
			update <- struct{}{}
		}
		if msg == command_main {
			nodeType = node_main
			update <- struct{}{}
		}
		if msg == command_backup {
			nodeType = node_backup
			update <- struct{}{}
		}
		if command, json := splitCommand(msg); command == command_updateNodeData {
			jsonStr = json
			update <- struct{}{}
		}
	})
}

func (register *Register) ToMain() {
	register.omipcClient.Notify(register.channel, command_main)
}

func (register *Register) ToBackup() {
	register.omipcClient.Notify(register.channel, command_backup)
}

func (register *Register) Close() {
	register.omipcClient.Notify(register.channel, command_close)
}

func (register *Register) Start() {
	register.omipcClient.Notify(register.channel, command_start)
}

func (register *Register) Stop() {
	register.omipcClient.Notify(register.channel, command_stop)
}

func (register *Register) UpdateData(data map[string]string) {
	jsonStr, _ := json.MarshalIndent(data, " ", "  ")
	register.omipcClient.Notify(register.channel, command_updateNodeData+":"+string(jsonStr))
}
