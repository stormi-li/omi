package omiweb

import (
	"log"
	"net/http"
	"strings"

	omiclient "github.com/stormi-li/omi/omi-manager"
)

type ProxyServer struct {
	router       *router
	omiWebClient *omiclient.Client
	serverName   string
	cache        *fileCache
}

func (proxyServer *ProxyServer) handleFunc(w http.ResponseWriter, r *http.Request) {
	domainNameResolution(r, proxyServer.router)
	httpProxy(w, r, proxyServer.cache)
	websocketProxy(w, r)
}

func (proxyServer *ProxyServer) SetCache(cacheDir string, maxSize int) {
	var err error
	proxyServer.cache, err = newFileCache(cacheDir, maxSize)
	if err != nil {
		panic(err)
	}
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
