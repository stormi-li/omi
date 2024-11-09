package main

import (
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi"
)

func requestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World send by http")
}

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	omi.NewServerClient(&redis.Options{Addr: redisAddr, Password: password}).NewRegister("helloworldserver", 1).Register("118.25.196.166:8081")
	http.HandleFunc("/request", requestHandler) // 注册 /request 路径的处理函数
	fmt.Println("Server is listening on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
