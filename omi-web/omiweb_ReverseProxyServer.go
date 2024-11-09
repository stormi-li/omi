package omiweb

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	omiclient "github.com/stormi-li/omi/omi-client"
)

type ReverseProxyServer struct {
	router       *Router
	omiWebClient *omiclient.Client
	serverName   string
	upgrader     websocket.Upgrader
}

func (proxy *ReverseProxyServer) handleFunc(w http.ResponseWriter, r *http.Request) {
	httpProxy(w, r, proxy.router, true)
	websocketProxy(w, r, proxy.router, true)
}

func (proxy *ReverseProxyServer) StartHttpProxy() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.handleFunc(w, r)
	})
	address := "127.0.0.1:80"
	proxy.omiWebClient.NewRegister("http反向代理", 1).Register(address)
	log.Println("omi web server: http反向代理 is running on http://" + address)
	err := http.ListenAndServe(":"+strings.Split(address, ":")[1], nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (proxy *ReverseProxyServer) StartHttpsProxy() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.handleFunc(w, r)
	})
	address := "127.0.0.1:80"
	proxy.omiWebClient.NewRegister(proxy.serverName, 1).Register(address)
	log.Println("omi web server: http反向代理 is running on http://" + address)
	err := http.ListenAndServe(":"+strings.Split(address, ":")[1], nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
