package omiweb

import (
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi"
)

type Client struct {
	router    *Router
	omiClient *omi.Client
}

func NewClient(redisClient *redis.Client, namespace string) *Client {
	omiClient := omi.NewClient(redisClient, namespace, omi.Server)
	return &Client{
		router:    newRouter(omiClient.NewSearcher()),
		omiClient: omiClient,
	}
}

func (omiweb *Client) NewManager() *Manager {
	return newManager(omiweb)
}

//go:embed src/*
var content embed.FS

func (omiweb *Client) Start(serverName, address string) {
	omiweb.start(serverName, address, getSourceFilePath()+"/src")
}

func (omiweb *Client) start(serverName, address string, srcPath string) {
	writeOmirequestPrefixToOmiJS(srcPath, address)
	copyResource(srcPath)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			data, err := content.ReadFile("src/index.html")
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

		r.URL.Path = "src" + r.URL.Path
		http.FileServer(http.FS(content)).ServeHTTP(w, r)
	})

	log.Println("omi web server: " + serverName + " is running on http://" + address)

	register := omiweb.omiClient.NewRegister(serverName, address)
	go register.StartOnMain(map[string]string{"omi web server": serverName})

	go http.ListenAndServe(":"+strings.Split(address, ":")[1], nil)
	<-register.CloseSignal
}

func writeOmirequestPrefixToOmiJS(src string, address string) {
	err := createOrReplaceFile(src, "omi.js", `export const omi_request_prefix = "http://`+address+`/omirequest";`)
	if err != nil {
		log.Fatalln(err)
	}
}

func copyResource(src string) {
	err := copyDir(src, "./src")
	if err != nil {
		log.Fatalln(err)
	}
}

func createOrReplaceFile(directory, filename, content string) error {
	// 组合路径
	fullPath := filepath.Join(directory, filename)

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); err == nil {
		// 文件存在，先删除
		if err := os.Remove(fullPath); err != nil {
			return fmt.Errorf("无法删除文件: %v", err)
		}
	} else if !os.IsNotExist(err) {
		// 如果其他错误（非文件不存在错误），返回错误
		return fmt.Errorf("检查文件时出错: %v", err)
	}

	// 创建或重新创建文件，并写入内容
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入文件时出错: %v", err)
	}

	return nil
}

func (omiweb *Client) forwardHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	// 以 '/' 分割路径，获取第一个参数
	parts := strings.Split(path, "/")

	address := omiweb.router.getAddress(parts[1])
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

func getSourceFilePath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return dir
}

func getStringAfterSecondSlash(input string) string {
	// 找到第一个 "/" 的位置
	firstSlashIndex := strings.Index(input, "/")
	if firstSlashIndex == -1 {
		return ""
	}

	// 在第一个 "/" 之后的字符串中再查找第二个 "/"
	secondSlashIndex := strings.Index(input[firstSlashIndex+1:], "/")
	if secondSlashIndex == -1 {
		return ""
	}

	// 计算第二个 "/" 的实际位置并截取子字符串
	secondSlashIndex += firstSlashIndex + 1
	return input[secondSlashIndex+1:]
}

// copyFile 复制单个文件
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// 使用 io.Copy 复制文件内容
	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

// copyDir 递归复制整个文件夹
func copyDir(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 构建目标路径
		relativePath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dstDir, relativePath)

		if info.IsDir() {
			// 创建目标文件夹
			return os.MkdirAll(dstPath, info.Mode())
		} else {
			// 复制文件
			return copyFile(path, dstPath)
		}
	})
}
