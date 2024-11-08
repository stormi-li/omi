package database

import (
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omi"
	"github.com/stormi-li/omi/app-server-login/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func init() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})

	addr, data := omi.NewConfigClient(redisClient, "omi-chat").NewSearcher().GetHighestPriorityServer("mysql")
	fmt.Println(data)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", data["username"], data["password"], addr, data["database"])
	fmt.Println(dsn)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// 自动迁移模型（创建表）
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
