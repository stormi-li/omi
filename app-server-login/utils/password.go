package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPasswordHash 检查密码是否正确
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var JWT_SECRET = []byte("omi")

func ValidateToken(tokenString string) (string, error) {
	// 解析并验证 Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证 Token 使用的签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return JWT_SECRET, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("Invalid or expired token")
	}

	// 提取声明中的用户信息
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["sub"].(string)
		return username, nil
	}

	return "", errors.New("Invalid token claims")
}

func GenerateJWT(username string) (string, error) {
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
	tokenString, err := token.SignedString(JWT_SECRET)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}