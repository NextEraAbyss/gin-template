package routes

import (
	"time"

	"gitee.com/NextEraAbyss/gin-template/internal/container"
	"gitee.com/NextEraAbyss/gin-template/middlewares"
	"gitee.com/NextEraAbyss/gin-template/models"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// SetupRoutes 配置所有路由
// 极简版本 - 只保留用户相关功能
func SetupRoutes(router *gin.Engine, db *gorm.DB, redisClient *redis.Client) {
	err := db.AutoMigrate(
		&models.User{}, // 只保留用户表
	)
	if err != nil {
		panic("数据库迁移失败: " + err.Error())
	}

	// 创建依赖注入容器
	newContainer := container.NewContainer(db, redisClient)
	newContainer.InitRepositories()
	newContainer.InitServices()
	newContainer.InitControllers()

	// Swagger 文档路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 保留必要的全局中间件
	router.Use(middlewares.RequestID())                  // 请求ID中间件，便于追踪请求
	router.Use(middlewares.Logger())                     // 日志中间件
	router.Use(middlewares.Recovery())                   // 恢复中间件
	router.Use(middlewares.ErrorHandler())               // 错误处理中间件
	router.Use(middlewares.CorsMiddleware())             // CORS中间件
	router.Use(middlewares.Security())                   // 安全中间件，添加安全相关HTTP头
	router.Use(middlewares.RateLimit(1000, time.Minute)) // 限流中间件，每分钟最多1000个请求

	// API路由组
	api := router.Group("/api/v1")

	// 用户相关路由
	users := api.Group("/users")
	users.GET("", newContainer.GetUserController().List)    // 获取用户列表
	users.GET("/:id", newContainer.GetUserController().Get) // 获取单个用户

	// 需要认证的用户路由
	userAuth := users.Group("")
	userAuth.Use(middlewares.AuthMiddleware())
	userAuth.PUT("/:id", newContainer.GetUserController().Update)                      // 更新用户信息
	userAuth.DELETE("/:id", newContainer.GetUserController().Delete)                   // 删除用户
	userAuth.POST("/change-password", newContainer.GetUserController().ChangePassword) // 修改密码
}

// RegisterRoutes 注册所有路由.
// 该函数用于注册不需要数据库连接的路由.
// 参数:
//   - router: Gin引擎实例.
func RegisterRoutes(router *gin.Engine) {
	// 这个函数可以保留为空，或者用于其他不需要数据库连接的路由注册.
}
