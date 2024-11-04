package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi/omi-manager/manager"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})
	address := "118.25.196.166:8080"

	managerMap := map[string]*manager.Manager{}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			// 读取 index.html 文件
			data, err := os.ReadFile(filepath.Join("src", "index.html"))
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
				managerMap[part[1]] = manager.NewManager(redisClient, part[1])
			}
			managerMap[part[1]].Handler(w, r)
			return
		}

		// 处理静态文件请求
		filePath := filepath.Join("src", r.URL.Path)
		http.ServeFile(w, r, filePath)
	})

	log.Println("omi web manager server is running on http://" + address)

	http.ListenAndServe(":"+strings.Split(address, ":")[1], nil)
}
