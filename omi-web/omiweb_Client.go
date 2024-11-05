package omiweb

import (
	"embed"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	omiclient "github.com/stormi-li/omi/omi_Client"
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
