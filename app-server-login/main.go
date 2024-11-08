package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi"
	"github.com/stormi-li/omi/app-server-login/controllers"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {

	client := omi.NewServerClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})
	register := client.NewRegister("用户登录注册服务", "118.25.196.166:7788")
	register.Start(1, map[string]string{})

	// 创建 Gin 路由
	r := gin.Default()

	// 注册路由
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/tokenValidation", controllers.TokeValidation)
	r.GET("/generate_qr", controllers.GenerateQRCode)
	r.GET("/qrValidation", controllers.LoginWebSocket)
	// 启动服务器
	r.Run(":7788")
}