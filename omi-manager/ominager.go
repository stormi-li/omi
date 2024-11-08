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

func Start(opts *redis.Options, address string) {

	managerMap := map[string]*Manager{}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // 使用 "*" 允许所有域
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

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
				managerMap[part[1]] = NewManager(opts)
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
