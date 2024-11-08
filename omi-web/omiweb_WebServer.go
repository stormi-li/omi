package omiweb

import (
	"embed"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/stormi-li/omi"
	omiclient "github.com/stormi-li/omi/omi-client"
)

type WebServer struct {
	router      *Router
	redisClient *redis.Client
	omiClient   *omiclient.Client
	serverName  string
	upgrader    websocket.Upgrader
	embedSource embed.FS
	embedModel  bool
}

func newWebServer(redisClient *redis.Client, serverName string) *WebServer {
	omiClient := omi.NewServerClient(redisClient.Options())
	return &WebServer{
		router:      newRouter(omiClient.NewSearcher()),
		redisClient: redisClient,
		omiClient:   omiClient,
		serverName:  serverName,
		upgrader:    websocket.Upgrader{},
	}
}

func (webServer *WebServer) SetEmbedSource(embedSource embed.FS) {
	webServer.embedSource = embedSource
	webServer.embedModel = true
}

func (webServer *WebServer) handleFunc(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, prefix_http_proxy) {
		httpProxy(w, r, webServer.router)
		return
	}

	if strings.HasPrefix(r.URL.Path, prefix_websocket_proxy) {
		websocketProxy(w, r, webServer.router)
		return
	}

	if r.URL.Path == "/" {
		var data []byte
		var err error
		if webServer.embedModel {
			data, err = webServer.embedSource.ReadFile(target_path + index_path)
		} else {
			data, err = os.ReadFile(target_path + index_path)
		}
		if err != nil {
			http.Error(w, "无法找到 "+index_path+" 文件", http.StatusNotFound)
			return
		}
		w.Write(data)
		return
	}

	if webServer.embedModel {
		r.URL.Path = target_path + r.URL.Path
		http.FileServer(http.FS(webServer.embedSource)).ServeHTTP(w, r)
	} else {
		http.ServeFile(w, r, target_path+r.URL.Path)
	}
}

func (webServer *WebServer) Listen(address string, weight int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		webServer.handleFunc(w, r)
	})
	omi.NewWebClient(webServer.redisClient.Options()).NewRegister(webServer.serverName, address).Start(weight, map[string]string{})
	log.Println("omi web server: " + webServer.serverName + " is running on http://" + address)
	err := http.ListenAndServe(":"+strings.Split(address, ":")[1], nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
