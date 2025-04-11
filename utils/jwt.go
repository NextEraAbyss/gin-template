package utils

import (
	"errors"
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
func InitJWTConfig(cfg *config.Config) {
	jwtSecret = []byte(cfg.JWT.Secret)
	jwtExpiration = time.Duration(cfg.JWT.ExpirationHours) * time.Hour
}

// Claims 自定义JWT声明
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(jwtExpiration)

	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  nowTime.Unix(),
			Issuer:    "gin-template",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString(jwtSecret)
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
