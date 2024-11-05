package ominager

import (
	"embed"
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
)

//go:embed src/*
var embedSource embed.FS

func Start(redisClient *redis.Client, address string) {

	managerMap := map[string]*Manager{}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			// 读取 index.html 文件
			data, err := embedSource.ReadFile("src/index.html")
			if err != nil {
				http.Error(w, "无法找到 index.html 文件", http.StatusNotFound)
				return
			}
			w.Write(data)
			return
		}

		// 如果请求的路径是特定的请求，转发处理
		part := strings.Split(r.URL.Path, "/")
		if len(part) > 1 && len(strings.Split(part[1], ".")) == 1 {
			if managerMap[part[1]] == nil {
				managerMap[part[1]] = NewManager(redisClient, part[1])
			}
			managerMap[part[1]].Handler(w, r)
			return
		}

		r.URL.Path = "src" + r.URL.Path
		http.FileServer(http.FS(embedSource)).ServeHTTP(w, r)
	})

	log.Println("omi web manager server is running on http://" + address)

	http.ListenAndServe(":"+strings.Split(address, ":")[1], nil)
}
