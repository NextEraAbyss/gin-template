package container

import (
	"gitee.com/NextEraAbyss/gin-template/controllers"
	"gitee.com/NextEraAbyss/gin-template/repositories"
	"gitee.com/NextEraAbyss/gin-template/services"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Container 依赖注入容器.
type Container struct {
	db           *gorm.DB
	redisClient  *redis.Client
	Repositories *Repositories
	Services     *Services
	Controllers  *Controllers
}

// Repositories 仓储层依赖.
type Repositories struct {
	User repositories.UserRepository
}

// Services 服务层依赖.
type Services struct {
	User services.UserService
}

// Controllers 控制器层依赖
type Controllers struct {
	User *controllers.UserController
}

// NewContainer 创建新的容器实例
func NewContainer(db *gorm.DB, redisClient *redis.Client) *Container {
	return &Container{
		db:          db,
		redisClient: redisClient,
	}
}

// InitRepositories 初始化仓储层
func (c *Container) InitRepositories() {
	c.Repositories = &Repositories{
		User: repositories.NewUserRepository(c.db),
	}
}

// InitServices 初始化服务层
func (c *Container) InitServices() {
	c.Services = &Services{
		User: services.NewUserService(c.Repositories.User),
	}
}

// InitControllers 初始化控制器层
func (c *Container) InitControllers() {
	c.Controllers = &Controllers{
		User: controllers.NewUserController(c.Services.User),
	}
}

// GetUserController 获取用户控制器
func (c *Container) GetUserController() *controllers.UserController {
	return c.Controllers.User
}
