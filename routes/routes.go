package routes

import (
	"gitee.com/NextEraAbyss/gin-template/controllers"
	"gitee.com/NextEraAbyss/gin-template/middlewares"
	"gitee.com/NextEraAbyss/gin-template/repositories"
	"gitee.com/NextEraAbyss/gin-template/services"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// SetupRoutes 配置所有路由
func SetupRoutes(router *gin.Engine, db *gorm.DB, redisClient *redis.Client) {
	// 创建依赖关系
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo, redisClient)
	userController := controllers.NewUserController(userService)

	// 公共路由
	router.POST("/login", userController.Login)

	// API路由组
	api := router.Group("/api")
	{
		// 用户路由
		users := api.Group("/users")
		{
			users.POST("", userController.Create)
			users.GET("", userController.GetAll)
			users.GET("/:id", userController.GetByID)

			// 需要认证的路由
			authUsers := users.Group("")
			authUsers.Use(middlewares.AuthMiddleware())
			{
				authUsers.PUT("/:id", userController.Update)
				authUsers.DELETE("/:id", userController.Delete)
			}
		}
	}
}
