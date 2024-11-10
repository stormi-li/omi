package omiweb

import (
	"log"
	"net/http"
	"strings"

	omiclient "github.com/stormi-li/omi/omi-client"
)

type ProxyServer struct {
	router       *Router
	omiWebClient *omiclient.Client
	serverName   string
}

func (proxyServer *ProxyServer) handleFunc(w http.ResponseWriter, r *http.Request) {
	domainNameResolution(r, proxyServer.router)
	httpProxy(w, r)
	websocketProxy(w, r)
}

func (proxyServer *ProxyServer) StartHttpProxy(address string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxyServer.handleFunc(w, r)
	})
	parts := strings.Split(address, ":")
	if len(parts) < 2 || parts[1] != "80" {
		panic("端口号必须为:80")
	}
	proxyServer.omiWebClient.NewRegister(proxyServer.serverName, 1).Register(address)
	log.Println("omi web server: " + proxyServer.serverName + " is running on http://" + address)
	err := http.ListenAndServe(":"+parts[1], nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
