package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stormi-li/omi"
)

// 设置 WebSocket 升级器
var upgrader = websocket.Upgrader{}

var connMap map[string]*websocket.Conn
var messageChan chan []byte

func init() {
	connMap = map[string]*websocket.Conn{}
	messageChan = make(chan []byte, 10000)
	go broadcastMessage()
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 升级 HTTP 请求为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()
	connMap[uuid.NewString()] = conn
	// 读取来自客户端的消息
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read message error:", err)
			break
		}
		messageChan <- p
	}
}

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func broadcastMessage() {
	for {
		msg := <-messageChan
		for _, conn := range connMap {
			conn.WriteMessage(1, msg)
		}
	}
}

func main() {
	address := "118.25.196.166:8181"
	// 设置 WebSocket 路由
	http.HandleFunc("/omichat", handleWebSocket)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})
	omiClient := omi.NewServerClient(redisClient, "omi-chat")
	omiClient.NewRegister("omi-chat-server", address).StartOnMain()
	// 启动 HTTP 服务器
	log.Println("omi-chat server started at ws://" + address)
	if err := http.ListenAndServe(":"+strings.Split(address, ":")[1], nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
