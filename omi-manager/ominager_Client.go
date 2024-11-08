package ominager

import (
	"embed"
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	opts *redis.Options
}

//go:embed src/*
var embedSource embed.FS

func (c *Client) Start(address string) {

	manager := NewManager(c.opts)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		part := strings.Split(r.URL.Path, "/")
		if r.URL.Path != "/" && len(part) > 1 && len(strings.Split(part[1], ".")) == 1 {
			manager.Handler(w, r)
			return
		}
		filePath := "src" + r.URL.Path
		if r.URL.Path == "/" {
			filePath = "src/index.html"
		}
		http.ServeFile(w, r, filePath)
		// http.FileServer(http.FS(embedSource)).ServeHTTP(w, r)
	})

	log.Println("omi web manager server is running on http://" + address)

	http.ListenAndServe(":"+strings.Split(address, ":")[1], nil)
}
