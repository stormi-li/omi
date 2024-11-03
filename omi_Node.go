package omi

import "encoding/json"

type Node struct {
	ServerType    ServerType
	ServerName    string
	State         string
	NodeType      string
	Address       string
	researdClient *Client
	searcher      *Searcher
	register      *Register
}

func newNode(serverType ServerType, serverName, state, nodeType, address string,
	researdClient *Client, searcher *Searcher) *Node {
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

func (node *Node) GetData() (map[string]string, string) {
	data := node.searcher.getData(node.ServerName, node.State, node.NodeType, node.Address)
	jsonByte, _ := json.MarshalIndent(data, " ", "  ")
	return data, string(jsonByte)
}

func (node *Node) ToMain() {
	node.NodeType = node_main
	node.register.ToMain()
}

func (node *Node) ToStandby() {
	node.NodeType = node_standby
	node.register.ToStandby()
}

func (node *Node) Close() {
	node.State = state_close
	node.register.Close()
}

func (node *Node) Start() {
	node.State = state_start
	node.register.Start()
}

func (node *Node) Stop() {
	node.State = state_stop
	node.register.Stop()
}

func (node *Node) UpdateData(data map[string]string) {
	node.register.UpdateData(data)
}

func (node *Node) ToString() string {
	bs, _ := json.MarshalIndent(node, " ", "  ")
	return string(bs)
}

func nodesToString(nodes []Node) string {
	bs, _ := json.MarshalIndent(nodes, " ", "  ")
	return string(bs)
}
