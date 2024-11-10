package omiweb

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

func pathRequestResolution(r *http.Request, router *Router) {
	serverName := strings.Split(r.URL.Path, "/")[1]
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/"+serverName)
	host := router.getAddress(serverName)
	r.URL.Host = host
}

func domainNameResolution(r *http.Request, router *Router) {
	host := router.getAddress(strings.Split(r.Host, ":")[0])
	r.URL.Host = host
}

func isWebSocketRequest(r *http.Request) bool {
	// 判断请求头中是否包含 WebSocket 升级请求特有的字段
	return r.Header.Get("Upgrade") == "websocket" &&
		r.Header.Get("Connection") == "Upgrade" &&
		r.Header.Get("Sec-WebSocket-Key") != ""
}

// 处理 HTTP 请求
func httpProxy(w http.ResponseWriter, r *http.Request) {
	if isWebSocketRequest(r) {
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   r.URL.Host,
	})
	proxy.ServeHTTP(w, r)
}

var upgrader = websocket.Upgrader{}

func websocketProxy(w http.ResponseWriter, r *http.Request) {
	if !isWebSocketRequest(r) {
		return
	}
	// 将客户端升级为WebSocket
	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}
	defer clientConn.Close()

	r.URL.Scheme = "ws"

	// 连接到后端WebSocket服务器
	targetConn, _, err := websocket.DefaultDialer.Dial(r.URL.String(), nil)
	if err != nil {
		log.Printf("无法连接到WebSocket服务器: %v", err)
		return
	}
	defer targetConn.Close()

	// 开始数据转发
	errChan := make(chan error, 2)

	go copyWebSocketData(targetConn, clientConn, errChan)
	go copyWebSocketData(clientConn, targetConn, errChan)

	// 等待传输结束
	<-errChan
}

// WebSocket数据复制
func copyWebSocketData(dst, src *websocket.Conn, errChan chan error) {
	for {
		msgType, msg, err := src.ReadMessage()
		if err != nil {
			errChan <- err
			return
		}
		if err := dst.WriteMessage(msgType, msg); err != nil {
			errChan <- err
			return
		}
	}
}
