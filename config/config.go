package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config 应用程序配置
type Config struct {
	Env    string
	Server struct {
		Host string
		Port string
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
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found. Using environment variables.")
	}

	config := &Config{}

	// 设置环境
	config.Env = getEnv("ENV", "development")

	// 服务器配置
	config.Server.Host = getEnv("SERVER_HOST", "localhost")
	config.Server.Port = getEnv("SERVER_PORT", "8080")

	// 数据库配置
	config.Database.Host = getEnv("DB_HOST", "localhost")
	config.Database.Port = getEnv("DB_PORT", "5432")
	config.Database.User = getEnv("DB_USER", "postgres")
	config.Database.Password = getEnv("DB_PASSWORD", "postgres")
	config.Database.Name = getEnv("DB_NAME", "gin_template")

	// Redis配置
	config.Redis.Host = getEnv("REDIS_HOST", "localhost")
	config.Redis.Port = getEnv("REDIS_PORT", "6379")
	config.Redis.Password = getEnv("REDIS_PASSWORD", "")
	config.Redis.DB = 0 // 默认使用DB 0

	// JWT配置
	config.JWT.Secret = getEnv("JWT_SECRET", "default-jwt-secret-never-use-this-in-production")
	// 默认过期时间为24小时
	config.JWT.ExpirationHours = 24

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
