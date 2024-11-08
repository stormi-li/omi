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
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})

	client := omi.NewServerClient(redisClient, "omi-chat")
	register := client.NewRegister("用户登录注册服务", "118.25.196.166:8084")
	register.StartOnMain()

	// 创建 Gin 路由
	r := gin.Default()

	// 注册路由
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/tokenValidation", controllers.TokeValidation)

	// 启动服务器
	r.Run(":8084")
}
