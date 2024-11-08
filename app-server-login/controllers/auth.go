package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stormi-li/omi/app-server-login/database"
	"github.com/stormi-li/omi/app-server-login/models"
	"github.com/stormi-li/omi/app-server-login/utils"
)

// 注册用户
func Register(c *gin.Context) {
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// 解析请求体
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// 检查用户是否已存在
	var existingUser models.User
	if err := database.DB.Where("username = ?", data.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(data.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 创建新用户
	user := models.User{
		Username: data.Username,
		Password: hashedPassword,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// 用户登录
func Login(c *gin.Context) {
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// 解析请求体
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// 查找用户
	var user models.User
	if err := database.DB.Where("username = ?", data.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 验证密码
	if !utils.CheckPasswordHash(data.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// 生成 JWT
	tokenString, err := generateJWT(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 登录成功，返回 JWT 给前端
	c.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"username": user.Username,
		"token":    tokenString,
	})
}

func TokeValidation(c *gin.Context) {
	// 从 Authorization 头部获取 Token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
		return
	}

	// 提取 Token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// 验证 Token
	username, err := utils.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// 如果 Token 有效，可以继续处理请求
	c.JSON(http.StatusOK, gin.H{
		"message":  "Access granted to protected endpoint",
		"username": username,
	})
}

func generateJWT(username string) (string, error) {
	// 设置 token 的过期时间，这里设置为 24 小时后
	expirationTime := time.Now().Add(24 * time.Hour)

	// 创建 JWT 的声明内容
	claims := &jwt.StandardClaims{
		Subject:   username,
		ExpiresAt: expirationTime.Unix(),
	}

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名 token
	tokenString, err := token.SignedString(utils.JWT_SECRET)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
