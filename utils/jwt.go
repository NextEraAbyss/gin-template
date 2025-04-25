package utils

import (
	"errors"
	"fmt"
	"time"

	"gitee.com/NextEraAbyss/gin-template/config"
	"github.com/golang-jwt/jwt/v4"
)

// 全局配置变量
var (
	jwtSecret     []byte
	jwtExpiration time.Duration
)

// InitJWTConfig 初始化JWT配置
func InitJWTConfig(config *config.Config) {
	jwtSecret = []byte(config.JWT.Secret)
	jwtExpiration = time.Duration(config.JWT.ExpirationHours) * time.Hour
}

// Claims 自定义JWT声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, username string) (string, error) {
	// 创建自定义的声明
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpiration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Subject:   fmt.Sprintf("%d", userID),
	}

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["username"] = username

	// 使用密钥签名 token
	return token.SignedString(jwtSecret)
}

// ParseToken 解析JWT令牌
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, errors.New("invalid token")
}
