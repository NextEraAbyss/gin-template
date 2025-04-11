package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gitee.com/NextEraAbyss/gin-template/internal/redis"
	"gitee.com/NextEraAbyss/gin-template/models"
	"gitee.com/NextEraAbyss/gin-template/repositories"
	"gitee.com/NextEraAbyss/gin-template/utils"
	redisClient "github.com/go-redis/redis/v8"
)

// 缓存键前缀和过期时间
const (
	UserCacheKeyPrefix = "user:"
	UserCacheDuration  = 10 * time.Minute
)

// UserService 接口定义了用户相关的业务逻辑
type UserService interface {
	CreateUser(user *models.User) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
	Authenticate(username, password string) (*models.User, error)
}

// userService 实现 UserService 接口
type userService struct {
	userRepo    repositories.UserRepository
	redisClient *redisClient.Client
}

// NewUserService 创建一个新的 UserService 实例
func NewUserService(userRepo repositories.UserRepository, redisClient *redisClient.Client) UserService {
	return &userService{
		userRepo:    userRepo,
		redisClient: redisClient,
	}
}

// getUserCacheKey 生成用户缓存键
func getUserCacheKey(id uint) string {
	return fmt.Sprintf("%s%d", UserCacheKeyPrefix, id)
}

// getUserByUsernameCacheKey 生成用户名缓存键
func getUserByUsernameCacheKey(username string) string {
	return fmt.Sprintf("%s%s", UserCacheKeyPrefix, username)
}

// CreateUser 创建新用户
func (s *userService) CreateUser(user *models.User) error {
	// 验证密码强度
	if err := utils.ValidatePasswordStrength(user.Password); err != nil {
		return err
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return s.userRepo.Create(user)
}

// GetUserByID 根据ID获取用户，优先从缓存获取
func (s *userService) GetUserByID(id uint) (*models.User, error) {
	// 尝试从缓存获取
	cacheKey := getUserCacheKey(id)
	cachedUser, err := redis.Get(cacheKey)

	// 如果找到缓存，反序列化并返回
	if err == nil {
		var user models.User
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			return &user, nil
		}
	}

	// 缓存未命中，从数据库获取
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 序列化并存入缓存
	userJSON, err := json.Marshal(user)
	if err == nil {
		redis.Set(cacheKey, userJSON, UserCacheDuration)
	}

	return user, nil
}

// GetUserByUsername 根据用户名获取用户，优先从缓存获取
func (s *userService) GetUserByUsername(username string) (*models.User, error) {
	// 尝试从缓存获取
	cacheKey := getUserByUsernameCacheKey(username)
	cachedUser, err := redis.Get(cacheKey)

	// 如果找到缓存，反序列化并返回
	if err == nil {
		var user models.User
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			return &user, nil
		}
	}

	// 缓存未命中，从数据库获取
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	// 序列化并存入缓存
	userJSON, err := json.Marshal(user)
	if err == nil {
		redis.Set(cacheKey, userJSON, UserCacheDuration)
	}

	return user, nil
}

// GetAllUsers 获取所有用户
func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAll()
}

// UpdateUser 更新用户信息
func (s *userService) UpdateUser(user *models.User) error {
	// 如果密码不为空，表示需要更新密码
	if user.Password != "" {
		// 验证密码强度
		if err := utils.ValidatePasswordStrength(user.Password); err != nil {
			return err
		}

		// 密码加密
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}

	err := s.userRepo.Update(user)
	if err != nil {
		return err
	}

	// 更新成功后，清除相关缓存
	idCacheKey := getUserCacheKey(user.ID)
	usernameCacheKey := getUserByUsernameCacheKey(user.Username)
	redis.Delete(idCacheKey)
	redis.Delete(usernameCacheKey)

	return nil
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(id uint) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 删除数据库记录
	err = s.userRepo.Delete(id)
	if err != nil {
		return err
	}

	// 删除成功后，清除相关缓存
	idCacheKey := getUserCacheKey(id)
	usernameCacheKey := getUserByUsernameCacheKey(user.Username)
	redis.Delete(idCacheKey)
	redis.Delete(usernameCacheKey)

	return nil
}

// Authenticate 用户认证
func (s *userService) Authenticate(username, password string) (*models.User, error) {
	user, err := s.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
