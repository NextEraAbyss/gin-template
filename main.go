package main

import (
	"path/filepath"

	"gitee.com/NextEraAbyss/gin-template/internal/database"
	"gitee.com/NextEraAbyss/gin-template/internal/redis"

	"gitee.com/NextEraAbyss/gin-template/config"
	"gitee.com/NextEraAbyss/gin-template/middlewares"
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

	// 应用中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.CorsMiddleware())
	router.Use(middlewares.ErrorHandlerMiddleware())
	utils.Debug("Middleware applied")

	// 注册路由
	routes.SetupRoutes(router, db, redisClient)
	utils.Info("Routes registered")

	// 启动服务器
	serverAddr := appConfig.Server.Host + ":" + appConfig.Server.Port
	utils.Info("Server running at %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		utils.Fatal("Failed to start server: %v", err)
	}
}
