package ominager

import omiclient "github.com/stormi-li/omi/omi-client"

func NewClient(serverSearcher *omiclient.Searcher, webSearcher *omiclient.Searcher, configSearcher *omiclient.Searcher) *Client {
	return &Client{
		serverSearcher: serverSearcher,
		webSearcher:    webSearcher,
		configSearcher: configSearcher,
	}
}
