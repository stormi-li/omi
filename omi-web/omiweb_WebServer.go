package omiweb

import (
	"embed"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	omiclient "github.com/stormi-li/omi/omi-client"
)

type WebServer struct {
	router      *Router
	redisClient *redis.Client
	omiClient   *omiclient.Client
	serverName  string
	weight      int
	upgrader    websocket.Upgrader
	embedSource embed.FS
	embedModel  bool
}

func newWebServer(redisClient *redis.Client, omiClient *omiclient.Client, serverName string, weight int) *WebServer {
	return &WebServer{
		router:      newRouter(omiClient.NewSearcher()),
		redisClient: redisClient,
		omiClient:   omiClient,
		serverName:  serverName,
		weight:      weight,
		upgrader:    websocket.Upgrader{},
	}
}

func (webServer *WebServer) EmbedSource(embedSource embed.FS) {
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

	filePath := r.URL.Path
	if r.URL.Path == "/" {
		filePath = index_path
	}
	filePath = target_path + filePath
	var data []byte
	if webServer.embedModel {
		data, _ = webServer.embedSource.ReadFile(filePath)
	} else {
		data, _ = os.ReadFile(filePath)
	}
	w.Write(data)
}

func (webServer *WebServer) Listen(address string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		webServer.handleFunc(w, r)
	})
	webServer.omiClient.NewRegister(webServer.serverName, webServer.weight).Register(address)
	log.Println("omi web server: " + webServer.serverName + " is running on http://" + address)
	err := http.ListenAndServe(":"+strings.Split(address, ":")[1], nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
