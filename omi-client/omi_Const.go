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
const namespace_separator = ":"
const const_expireTime = 2 * time.Second
const const_listenWaitTime = 1 * time.Second

const Prefix_Config = "stormi:config:"
const Prefix_Server = "stormi:server:"
const Prefix_Web = "stormi:web:"

var Server = "Server"
var Config = "Config"
var Web = "Web"
