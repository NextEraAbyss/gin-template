package utils

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 对密码进行哈希处理
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash 比较密码与哈希值
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePasswordStrength 验证密码强度
// 密码必须同时包含大小写字母、数字和特殊字符，且长度至少为8位
func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasUpperCase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLowerCase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasSpecialChar := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	if !hasNumber || !hasUpperCase || !hasLowerCase || !hasSpecialChar {
		return errors.New("password must contain at least one number, one uppercase letter, one lowercase letter, and one special character")
	}

	return nil
}
