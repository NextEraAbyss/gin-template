package routes

import (
	"time"

	"gitee.com/NextEraAbyss/gin-template/internal/container"
	"gitee.com/NextEraAbyss/gin-template/middleware"
	"gitee.com/NextEraAbyss/gin-template/models"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// SetupRoutes 配置所有路由.
// 该函数负责设置应用的所有路由，包括API路由和中间件.
// 参数:
//   - router: Gin引擎实例.
//   - db: 数据库连接.
//   - redisClient: Redis客户端.
func SetupRoutes(router *gin.Engine, db *gorm.DB, redisClient *redis.Client) {
	err := db.AutoMigrate(
		&models.User{},    // 用户表.
		&models.Article{}, // 文章表.
	)
	if err != nil {
		panic("数据库迁移失败: " + err.Error())
	}

	// 创建依赖注入容器.
	newContainer := container.NewContainer(db, redisClient)
	newContainer.InitRepositories()
	newContainer.InitServices()
	newContainer.InitControllers()

	// Swagger 文档路由 - 在全局中间件之前注册，避免 CSP 限制.
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 全局中间件.
	router.Use(middleware.RequestID())      // 请求ID中间件.
	router.Use(middleware.Logger())         // 日志中间件.
	router.Use(middleware.Recovery())       // 恢复中间件.
	router.Use(middleware.ErrorHandler())   // 错误处理中间件.
	router.Use(middleware.CorsMiddleware()) // CORS中间件.
	router.Use(middleware.Security())
	router.Use(middleware.RateLimit(100, time.Minute)) // 限制每分钟100个请求.

	// API路由组.
	api := router.Group("/api/v1")

	// 用户认证相关路由.
	auth := api.Group("/auth")
	auth.POST("/login", newContainer.GetAuthController().Login)
	auth.POST("/register", newContainer.GetAuthController().Register)

	// 用户相关路由.
	users := api.Group("/users")
	// 公开接口 - 不需要认证.
	users.GET("", newContainer.GetUserController().List)    // 获取所有用户.
	users.GET("/:id", newContainer.GetUserController().Get) // 根据ID获取用户.

	// 需要认证的路由组.
	authUsers := users.Group("")
	authUsers.Use(middleware.AuthMiddleware())
	authUsers.PUT("/:id", newContainer.GetUserController().Update)                      // 更新用户信息.
	authUsers.DELETE("/:id", newContainer.GetUserController().Delete)                   // 删除用户.
	authUsers.POST("/change-password", newContainer.GetUserController().ChangePassword) // 修改密码.

	// 文章相关路由.
	articles := api.Group("/articles")
	// 公开接口 - 不需要认证.
	articles.GET("", newContainer.GetArticleController().List)    // 获取文章列表.
	articles.GET("/:id", newContainer.GetArticleController().Get) // 获取单个文章.

	// 需要认证的接口.
	authArticles := articles.Group("")
	authArticles.Use(middleware.AuthMiddleware())
	authArticles.POST("", newContainer.GetArticleController().Create)       // 创建文章.
	authArticles.PUT("/:id", newContainer.GetArticleController().Update)    // 更新文章.
	authArticles.DELETE("/:id", newContainer.GetArticleController().Delete) // 删除文章.
}

// RegisterRoutes 注册所有路由.
// 该函数用于注册不需要数据库连接的路由.
// 参数:
//   - router: Gin引擎实例.
func RegisterRoutes(router *gin.Engine) {
	// 这个函数可以保留为空，或者用于其他不需要数据库连接的路由注册.
}
