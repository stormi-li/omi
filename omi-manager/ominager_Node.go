package ominager

import (
	"encoding/json"

	omiclient "github.com/stormi-li/omi/omi_Client"
)

type Node struct {
	ServerType string
	ServerName string
	State      string
	NodeType   string
	Address    string
	omiClient  *omiclient.Client
	searcher   *omiclient.Searcher
	register   *omiclient.Register
}

func newNode(serverType string, serverName, state, nodeType, address string, omiClient *omiclient.Client, searcher *omiclient.Searcher) *Node {
	register := omiClient.NewRegister(serverName, address)
	return &Node{
		ServerType: serverType,
		ServerName: serverName,
		State:      state,
		NodeType:   nodeType,
		Address:    address,
		omiClient:  omiClient,
		searcher:   searcher,
		register:   register,
	}
}

func (node *Node) ToMain() {
	node.register.ToMain()
}

func (node *Node) ToBackup() {
	node.register.ToBackup()
}

func (node *Node) Start() {
	node.register.ToStart()
}

func (node *Node) Stop() {
	node.register.ToStop()
}

func (node *Node) ToString() string {
	bs, _ := json.MarshalIndent(node, " ", "  ")
	return string(bs)
}

func nodesToString(nodes []Node) string {
	bs, _ := json.MarshalIndent(nodes, " ", "  ")
	return string(bs)
}
