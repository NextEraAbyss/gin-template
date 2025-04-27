package utils

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"gitee.com/NextEraAbyss/gin-template/config"
	"github.com/golang-jwt/jwt/v4"
)

// 全局配置变量
var (
	jwtSecret     []byte
	jwtExpiration time.Duration
)

// 常见错误
var (
	ErrInvalidToken = errors.New("invalid token")
)

// InitJWTConfig 初始化JWT配置
func InitJWTConfig(config *config.Config) {
	jwtSecret = []byte(config.JWT.Secret)
	jwtExpiration = time.Duration(config.JWT.ExpirationHours) * time.Hour
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, username string) (string, error) {
	// 创建标准声明
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpiration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Subject:   fmt.Sprintf("%d", userID),
		Issuer:    username, // 使用Issuer存储username
	}

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名 token
	return token.SignedString(jwtSecret)
}

// ParseToken 解析JWT令牌
func ParseToken(token string) (*jwt.RegisteredClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*jwt.RegisteredClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, ErrInvalidToken
}

// GetUserIDFromToken 从令牌中提取用户ID
func GetUserIDFromToken(tokenString string) (uint, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return 0, err
	}

	userID, err := strconv.ParseUint(claims.Subject, 10, 32)
	if err != nil {
		return 0, errors.New("invalid token subject")
	}

	return uint(userID), nil
}
