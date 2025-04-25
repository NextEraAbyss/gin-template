package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 对密码进行加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword 验证密码是否正确
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateRandomPassword 生成随机密码
func GenerateRandomPassword() (string, error) {
	// 生成32字节的随机数
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random password: %w", err)
	}
	// 使用base64编码，并去掉特殊字符
	return base64.URLEncoding.EncodeToString(b)[:12], nil
}

// ValidatePasswordStrength 验证密码强度
// 密码必须同时包含大小写字母、数字和特殊字符，且长度至少为8位
func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hasNumber := regexp.MustCompile(`\d`).MatchString(password)
	hasUpperCase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLowerCase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasSpecialChar := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]`).MatchString(password)

	if !hasNumber || !hasUpperCase || !hasLowerCase || !hasSpecialChar {
		return errors.New("password must contain at least one number, one uppercase letter, one lowercase letter, and one special character")
	}

	return nil
}

// 生成随机盐值
func GenerateSalt() (string, error) {
	// 生成32字节的随机数
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	// 使用base64编码，并去掉特殊字符
	return base64.URLEncoding.EncodeToString(b)[:12], nil
}

// GenerateRandomString 生成指定长度的随机字符串.
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return string(bytes), nil
}
