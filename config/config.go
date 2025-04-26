package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用程序配置.
type Config struct {
	Env    string
	Server struct {
		Host string
		Port int
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	Redis struct {
		Host     string
		Port     string
		Password string
		DB       int
	}
	// 添加JWT配置
	JWT struct {
		Secret          string
		ExpirationHours int
	}
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 直接加载项目根目录的.env文件
	err := godotenv.Load()
	if err != nil {
		log.Println("警告: 未找到.env文件，将使用默认值或环境变量")
	} else {
		log.Println("配置已加载: .env (根目录)")
	}

	config := &Config{}

	// 设置环境
	config.Env = getEnv("ENV", "development")

	// 服务器配置
	config.Server.Host = getEnv("SERVER_HOST", "0.0.0.0")
	port, _ := strconv.Atoi(getEnv("SERVER_PORT", "9999"))
	config.Server.Port = port

	// 数据库配置
	config.Database.Host = getEnv("DB_HOST", "localhost")
	config.Database.Port = getEnv("DB_PORT", "3306")
	config.Database.User = getEnv("DB_USER", "root")
	config.Database.Password = getEnv("DB_PASSWORD", "")
	config.Database.Name = getEnv("DB_NAME", "gin_template")

	// Redis配置
	config.Redis.Host = getEnv("REDIS_HOST", "localhost")
	config.Redis.Port = getEnv("REDIS_PORT", "6379")
	config.Redis.Password = getEnv("REDIS_PASSWORD", "")
	config.Redis.DB = 0 // 默认使用DB 0

	// JWT配置
	config.JWT.Secret = getEnv("JWT_SECRET", "default-jwt-secret-never-use-this-in-production")
	expirationHours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		expirationHours = 24 // 默认24小时
	}
	config.JWT.ExpirationHours = expirationHours

	// 打印当前使用的配置信息
	printConfig(config)

	return config
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// printConfig 打印配置信息（隐藏敏感信息）
func printConfig(config *Config) {
	if config.Env != "production" {
		fmt.Println("=== 应用配置信息 ===")
		fmt.Printf("环境: %s\n", config.Env)
		fmt.Printf("服务器: %s:%d\n", config.Server.Host, config.Server.Port)
		fmt.Printf("数据库: %s:%s/%s\n", config.Database.Host, config.Database.Port, config.Database.Name)
		fmt.Printf("Redis: %s:%s\n", config.Redis.Host, config.Redis.Port)
		fmt.Println("===================")
	}
}
