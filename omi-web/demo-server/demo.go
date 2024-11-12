package main

import (
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	serverManager := omi.NewServerManager(&redis.Options{Addr: redisAddr, Password: password})
	register := serverManager.NewRegister("hello_server", 1)
	register.Register("118.25.196.166:8081")

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello", r.URL.Query().Get("name"), ", welcome to use omi")
	})
	http.ListenAndServe(":8081", nil)
}
