package omiweb

import (
	"embed"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	omiclient "github.com/stormi-li/omi/omi-client"
)

type Client struct {
	router      *Router
	redisClient *redis.Client
	omiClient   *omiclient.Client
	serverName  string
	namespace   string
	address     string
}

func (omiweb *Client) GenerateTemplate() {
	copyResource(getSourceFilePath() + "/TemplateSource")
}

func (omiweb *Client) Live() {
	omiweb.start()
}

func (omiweb *Client) Start(embedSource embed.FS) {
	omiweb.start(embedSource)
}

func (omiweb *Client) start(embedSources ...embed.FS) {
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
			omiweb.forwardHandler(w, r)
			return
		}

		if embedModel {
			r.URL.Path = "src" + r.URL.Path
			http.FileServer(http.FS(embedSource)).ServeHTTP(w, r)
		} else {
			http.ServeFile(w, r, "src/"+r.URL.Path)
		}
	})

	log.Println("omi web server: " + omiweb.serverName + " is running on http://" + omiweb.address)

	register := omiweb.omiClient.NewRegister(omiweb.serverName, omiweb.address)
	go register.StartOnMain(map[string]string{"omi web server": omiweb.serverName})

	http.ListenAndServe(":"+strings.Split(omiweb.address, ":")[1], nil)
}

func (omiweb *Client) forwardHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	// 以 '/' 分割路径，获取第一个参数
	parts := strings.Split(path, "/")

	address := omiweb.router.getAddress(parts[1])

	//未获取到地址
	if address == "" {
		http.Error(w, "未获取到地址:"+parts[1], http.StatusInternalServerError)
		return
	}

	targetURL := address + "/" + getStringAfterSecondSlash(path)
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
