package omiweb

import (
	"embed"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/stormi-li/omi"
	omiclient "github.com/stormi-li/omi/omi-client"
)

type Client struct {
	router      *Router
	redisClient *redis.Client
	omiClient   *omiclient.Client
	serverName  string
	namespace   string
	address     string
	upgrader    websocket.Upgrader
}

func (omiweb *Client) GenerateTemplate() {
	copyResource(getSourceFilePath() + "/TemplateSource")
}



func (omiweb *Client) Listen(address string, embedSources ...embed.FS) {
	omiweb.address = address
	var embedSource embed.FS
	embedModel := false
	if len(embedSources) > 0 {
		embedSource = embedSources[0]
		embedModel = true
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			var data []byte
			var err error
			if embedModel {
				data, err = embedSource.ReadFile("src/index.html")
			} else {
				data, err = os.ReadFile("src/index.html")
			}
			if err != nil {
				http.Error(w, "无法找到 index.html 文件", http.StatusNotFound)
				return
			}
			w.Write(data)
			return
		}

		part := strings.Split(r.URL.Path, "/")

		if len(part) > 1 && part[1] == const_omirequest {
			omiweb.requestForwardHandler(w, r)
			return
		}
		if len(part) > 1 && part[1] == const_omiwebsocket {
			omiweb.websocketForwardHandler(w, r)
			return
		}

		if embedModel {
			r.URL.Path = "src" + r.URL.Path
			http.FileServer(http.FS(embedSource)).ServeHTTP(w, r)
		} else {
			http.ServeFile(w, r, "src/"+r.URL.Path)
		}
	})
	omi.NewServerClient(omiweb.redisClient, omiweb.namespace).NewRegister("web-server", address).StartOnMain()
	log.Println("omi web server: " + omiweb.serverName + " is running on http://" + omiweb.address)
	http.ListenAndServe(":"+strings.Split(omiweb.address, ":")[1], nil)
}

func (omiweb *Client) getTargetURL(r *http.Request) string {
	path := strings.TrimPrefix(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	path = getStringAfterSecondSlash(path)

	if len(parts) < 2 {
		return ""
	}
	address := omiweb.router.getAddress(parts[1])
	if address == "" {
		return ""
	}
	return address + "/" + path
}

func (omiweb *Client) requestForwardHandler(w http.ResponseWriter, r *http.Request) {
	targetURL := omiweb.getTargetURL(r)
	if targetURL == "" {
		http.Error(w, "请求地址不存在或请求非法", http.StatusInternalServerError)
		return
	}
	// 创建一个 HTTP 请求，将 A 发送给 B 的请求原样转发给 C
	req, err := http.NewRequest(r.Method, "http://"+targetURL, r.Body)
	if err != nil {
		http.Error(w, "无法创建请求", http.StatusInternalServerError)
		return
	}
	// 复制请求头，以保持请求的原始头信息
	req.Header = r.Header

	// 使用 HTTP 客户端发送请求到 C
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "请求转发失败", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 将 C 的响应头写回给 A
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// 设置返回状态码为 C 返回的状态码
	w.WriteHeader(resp.StatusCode)

	// 将 C 的响应体原封不动地返回给 A
	io.Copy(w, resp.Body)
}

func (omiweb *Client) websocketForwardHandler(w http.ResponseWriter, r *http.Request) {
	targetURL := omiweb.getTargetURL(r)

	if targetURL == "" {
		http.Error(w, "请求地址不存在或请求非法", http.StatusInternalServerError)
		return
	}
	// 与 A 建立 WebSocket 连接
	clientConn, err := omiweb.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("与前端建立websocket连接失败:", err)
		return
	}
	defer clientConn.Close()

	// 与 C 建立 WebSocket 连接
	cConn, _, err := websocket.DefaultDialer.Dial("ws://"+targetURL, nil)
	if err != nil {
		log.Println("与服务端建立websocket连接失败:", err)
		return
	}
	defer cConn.Close()

	close := make(chan struct{}, 1)
	// 将 A 发来的消息转发给 C
	go forwardToC(clientConn, cConn, close)

	// 将 C 发来的消息转发回 A
	go forwardToA(cConn, clientConn, close)

	// 阻塞主协程直到连接关闭
	<-close
}

// 转发 A 发来的消息给 C
func forwardToC(aConn, cConn *websocket.Conn, close chan struct{}) {
	for {
		_, message, err := aConn.ReadMessage()
		if err != nil {
			close <- struct{}{}
			return
		}

		// 将消息转发给 C
		err = cConn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			close <- struct{}{}
			return
		}
	}
}

// 转发 C 发来的消息给 A
func forwardToA(cConn, aConn *websocket.Conn, close chan struct{}) {
	for {
		_, message, err := cConn.ReadMessage()
		if err != nil {
			close <- struct{}{}
			return
		}

		// 将消息转发给 A
		err = aConn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			close <- struct{}{}
			return
		}
	}
}
