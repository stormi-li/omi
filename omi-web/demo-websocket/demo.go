package main

import (
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/stormi-li/omi"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 升级 HTTP 请求到 WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	// 向客户端发送 "Hello World" 消息
	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello World send by websocket"))
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
}

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	omi.NewServerManager(&redis.Options{Addr: redisAddr, Password: password}).NewRegister("helloworldwebsocket", 1).Register("118.25.196.166:8082")

	http.HandleFunc("/request", wsHandler) // 将 /request 路径映射到 wsHandler
	fmt.Println("WebSocket server listening on :8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
