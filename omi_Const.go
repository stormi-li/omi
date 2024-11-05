package omi

import "time"

const command_updateNodeData = "updateNodeData"
const command_start = "start"
const command_stop = "stop"
const state_start = "start"
const state_stop = "stop"
const const_configPrefix = "stormi:config:"
const const_serverPrefix = "stormi:server:"
const const_mqPrefix = "stormi:mq:"
const const_separator = ":"
const command_toBackup = "backup"
const command_toMain = "main"
const node_backup = "backup"
const node_main = "main"

const const_retryWaitTime = 500 * time.Millisecond
const const_maxRetryCount = 10
const const_expireTime = 2 * time.Second

const const_listenWaitTime = 1*time.Second

type ServerType string

var Server ServerType = "Server"
var MQ ServerType = "MQ"
var Config ServerType = "Config"
