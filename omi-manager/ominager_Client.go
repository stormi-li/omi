package ominager

import (
	"embed"
	"log"
	"net/http"
	"os"
	"strings"

	omiclient "github.com/stormi-li/omi/omi-client"
)

type Client struct {
	serverSearcher *omiclient.Searcher
	webSearcher    *omiclient.Searcher
	configSearcher *omiclient.Searcher
}

//go:embed src/*
var embedSource embed.FS

func (c *Client) Start(address string) {

	manager := NewManager(c.serverSearcher, c.webSearcher, c.configSearcher)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		part := strings.Split(r.URL.Path, "/")
		if r.URL.Path != "/" && len(part) > 1 && len(strings.Split(part[1], ".")) == 1 {
			manager.Handler(w, r)
			return
		}

		filePath := r.URL.Path
		if r.URL.Path == "/" {
			filePath = "src/index.html"
		}
		filePath = "src" + filePath
		var data []byte
		// data, _ = embedSource.ReadFile(filePath)
		data, _ = os.ReadFile(filePath)
		w.Write(data)
	})

	log.Println("omi web manager server is running on http://" + address)

	http.ListenAndServe(":"+strings.Split(address, ":")[1], nil)
}
