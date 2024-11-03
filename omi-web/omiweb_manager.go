package omiweb

type Manager struct {
	omiweb *Client
}

func newManager(Omiweb *Client) *Manager {
	return &Manager{
		omiweb: Omiweb,
	}
}

func (manager *Manager) Start(serverName, address string) {
	manager.omiweb.start(serverName, address, getSourceFilePath()+"/webManagerSrc")
}
