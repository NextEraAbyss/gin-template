package main

import (
	"fmt"
	"path/filepath"
	"strconv"

	"gitee.com/NextEraAbyss/gin-template/internal/database"
	"gitee.com/NextEraAbyss/gin-template/internal/redis"

	"gitee.com/NextEraAbyss/gin-template/config"
	"gitee.com/NextEraAbyss/gin-template/routes"
	"gitee.com/NextEraAbyss/gin-template/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	appConfig := config.LoadConfig()

	// 初始化日志系统
	logDir := filepath.Join(".", "logs")
	utils.InitLogFile(logDir)
	if appConfig.Env == "production" {
		utils.InitLogger(utils.INFO, nil, nil, false)
	} else {
		utils.InitLogger(utils.DEBUG, nil, nil, true)
	}

	utils.Info("Starting application in %s mode", appConfig.Env)

	// 初始化JWT配置
	utils.InitJWTConfig(appConfig)
	utils.Debug("JWT configuration initialized")

	// 设置gin模式
	if appConfig.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 初始化数据库
	db := database.Init(appConfig)
	utils.Info("Database initialized")

	// 初始化Redis
	redisClient := redis.Init(appConfig)
	utils.Info("Redis initialized")

	// 初始化路由
	router := gin.New()

	// 注册路由
	routes.SetupRoutes(router, db, redisClient)
	utils.Info("Routes registered")

	// 启动服务器
	port, err := strconv.Atoi(appConfig.Server.Port)
	if err != nil {
		utils.Error("Invalid port number: %v", err)
		return
	}
	utils.Info("Server is running on port %d", port)
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		utils.Error("Server failed to start: %v", err)
	}
}
