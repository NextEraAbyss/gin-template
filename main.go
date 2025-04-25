package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitee.com/NextEraAbyss/gin-template/config"
	"gitee.com/NextEraAbyss/gin-template/internal/database"
	"gitee.com/NextEraAbyss/gin-template/routes"
	"gitee.com/NextEraAbyss/gin-template/utils"

	_ "gitee.com/NextEraAbyss/gin-template/docs" // 导入生成的文档.
	"github.com/gin-gonic/gin"
)

const (
	EnvProduction = "production"
)

// @title           Gin API Template
// @version         1.0
// @description     This is a sample server for a Gin API template.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9999
// @BasePath  /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// 加载配置.
	cfg := config.LoadConfig()

	// 初始化日志.
	utils.InitLogger(utils.INFO, os.Stdout, os.Stderr, true)

	// 初始化数据库.
	db := database.Init(cfg)

	// 初始化Redis.
	redisClient := database.InitRedis(cfg)

	// 设置Gin模式.
	if cfg.Env == EnvProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin引擎.
	router := gin.Default()

	// 设置路由.
	routes.SetupRoutes(router, db, redisClient)

	// 创建HTTP服务器.
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// 在goroutine中启动服务器.
	go func() {
		utils.Infof("Server is running on port %d", cfg.Server.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.Infof("Shutting down server...")

	// 设置关闭超时.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器.
	if err := srv.Shutdown(ctx); err != nil {
		utils.Fatalf("Server forced to shutdown: %v", err)
	}

	utils.Infof("Server exiting")
}
