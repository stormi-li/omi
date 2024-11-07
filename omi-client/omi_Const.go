package omiclient

import "time"

const command_start = "start"
const command_stop = "stop"
const state_start = "start"
const state_stop = "stop"
const command_toBackup = "backup"
const command_toMain = "main"
const node_backup = "backup"
const node_main = "main"

const const_retryWaitTime = 500 * time.Millisecond
const const_maxRetryCount = 10
const const_expireTime = 2 * time.Second

const const_listenWaitTime = 1 * time.Second

const Prefix_Config = "stormi:config:"
const Prefix_Server = "stormi:server:"
const Prefix_MQ = "stormi:mq:"
const Prefix_Web = "stormi:web:"

var Server = "Server"
var MQ = "MQ"
var Config = "Config"
var Web = "Config"
