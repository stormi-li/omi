package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/stormi-li/omi"
)

// 设置 WebSocket 升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 在此可以添加允许的 Origin 校验
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 升级 HTTP 请求为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// 读取来自客户端的消息
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read message error:", err)
			break
		}

		// 输出收到的消息
		fmt.Println("Received message:", string(p))

		go func() {
			for i := 0; i < 1000; i++ {
				// 回复消息
				if err := conn.WriteMessage(messageType, []byte(strconv.Itoa(i))); err != nil {
					log.Println("Write message error:", err)
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
		}()
	}
}

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	// 设置 WebSocket 路由
	http.HandleFunc("/ws", handleWebSocket)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})
	omiClient := omi.NewServerClient(redisClient, "omi-namespace")
	go omiClient.NewRegister("websocketserver", "118.25.196.166:8181").StartOnMain()
	// 启动 HTTP 服务器
	log.Println("WebSocket server started at ws://localhost:8181")
	if err := http.ListenAndServe(":8181", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
