package ominager

import (
	"encoding/json"

	"github.com/stormi-li/omi"
)

type Node struct {
	ServerType    omi.ServerType
	ServerName    string
	State         string
	NodeType      string
	Address       string
	researdClient *omi.Client
	searcher      *omi.Searcher
	register      *omi.Register
}

func newNode(serverType omi.ServerType, serverName, state, nodeType, address string,
	researdClient *omi.Client, searcher *omi.Searcher) *Node {
	register := researdClient.NewRegister(serverName, address)
	return &Node{
		ServerType:    serverType,
		ServerName:    serverName,
		State:         state,
		NodeType:      nodeType,
		Address:       address,
		researdClient: researdClient,
		searcher:      searcher,
		register:      register,
	}
}

func (node *Node) ToMain() {
	node.NodeType = nodeType_main
	node.register.ToMain()
}

func (node *Node) ToBackup() {
	node.NodeType = nodeType_backup
	node.register.ToBackup()
}

func (node *Node) Start() {
	node.State = state_start
	node.register.Start()
}

func (node *Node) Stop() {
	node.State = state_stop
	node.register.Stop()
}

func (node *Node) ToString() string {
	bs, _ := json.MarshalIndent(node, " ", "  ")
	return string(bs)
}

func nodesToString(nodes []Node) string {
	bs, _ := json.MarshalIndent(nodes, " ", "  ")
	return string(bs)
}
