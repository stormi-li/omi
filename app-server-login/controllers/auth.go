package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/skip2/go-qrcode"
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
	tokenString, err := utils.GenerateJWT(user.Username)
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

// token验证
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

var (
	sessions     = make(map[string]string) // sessionID to token mapping
	sessionMutex = &sync.Mutex{}           // protect sessions map
	upgrader     = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
)

// 生成二维码
func GenerateQRCode(c *gin.Context) {
	sessionID := uuid.NewString() // 自定义函数生成唯一 sessionID
	sessionMutex.Lock()
	sessions[sessionID] = "" // 初始化 sessionID 状态
	sessionMutex.Unlock()

	// 生成二维码 URL
	qrCodeURL := "http://118.25.196.166:8084/scan?session_id=" + sessionID
	png, _ := qrcode.Encode(qrCodeURL, qrcode.Medium, 256)
	c.JSON(http.StatusOK, gin.H{"session_id": sessionID, "qr_code": png})
}

type SessionRequest struct {
	SessionID string `json:"session_id"`
}

// WebSocket 连接，用于检查登录状态
func LoginWebSocket(c *gin.Context) {
	// 升级 HTTP 连接为 WebSocket 连接
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}
	defer ws.Close()

	var sessionID string

	// 读取并处理来自客户端的消息
	_, message, err := ws.ReadMessage()
	if err != nil {
		return
	}

	var req SessionRequest
	if err := json.Unmarshal(message, &req); err != nil {
		return
	}
	sessionID = req.SessionID
	sessions[sessionID] = c.Query("sessionID")
	for {
		// 检查 sessionID 对应的 token
		sessionMutex.Lock()
		token, exists := sessions[sessionID]
		sessionMutex.Unlock()

		// 如果 token 存在，则发送并关闭 WebSocket 连接
		if exists && token != "" {
			err := ws.WriteMessage(websocket.TextMessage, []byte(token))
			delete(sessions, sessionID)
			if err != nil {
				return
			}
			break
		}
		// 如果 token 不存在，等待 1 秒再检查
		time.Sleep(1 * time.Second)
	}
}

func ConfirmLogin(c *gin.Context) {
	sessionID := c.Query("session_id")
	token := "generated_user_token" // 实际生产环境中生成用户 Token

	sessionMutex.Lock()
	sessions[sessionID] = token
	sessionMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Login confirmed"})
}
